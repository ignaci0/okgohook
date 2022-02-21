// This package primaraly's purpose is to provide a router for implementing a
// webhook and conviniently handle the intent's fulfillment requests towards the webhook.
//
// The Router implements' net/http Handler interface to serve Google Actions
// webhook requests that are provided in the routes. Creating a webhook is as simple
// as:
//
//	router := okgohook.NewRouter()
//	router.HandleIntent("HelloWorld", func(aog.FulfillmentRequest) *aog.Fulfillmentresponse { ... })
//	http.Handle(router) // Other routers such as gorilla/mux can be used as well
//	log.Fatal(http.ListenAndServe(":6060", nil))
//
// The default router shall also respond the Google Health checks sent for
// analytics
package okgohook

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ignaci0/okgohook/aog"
)

// Router implements net/http's Handler interface to serve POST requests
// from Google Actions by serving the configured routes.
type Router struct {
	routes []*Route
	aud    string
}

var initialize sync.Once
var google_jwks_uri string = "https://www.googleapis.com/oauth2/v1/certs"

//var google_jwks_uri string = "https://www.googleapis.com/oauth2/v3/certs"
var google_keys map[string]string

const ACTIONS_INTENT_HEALTH_CHECK = "actions.intent.HEALTH_CHECK"

// NewRouter allocates and initializes a pointer to a Router struct adding also the
// default intent handlers.
//	var router *okgohook.Router
//	router = okgohook.NewRouter()
func NewRouter() *Router {
	return &Router{
		routes: []*Route{
			&Route{
				function: handleDefaultHealthCheck,
				intent:   ACTIONS_INTENT_HEALTH_CHECK,
			},
		},
	}
}

func (this *Router) Proxy(url string) {
}

func (this *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		//Return invalid request
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Authorization
	if !this.authorized(req.Header.Get("google-assistant-signature")) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&map[string]string{"error": "Unauthorized"})
		return
	}

	var r aog.FulfillmentRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {

		//Return malformed request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Match request to the handlers...
	for _, route := range this.routes {
		if route.isMatch(&r) {
			resp := route.function(&r)
			if resp == nil {
				//If something went wrong in the handler, let's give a chance to the next
				continue
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"status": { "code": 404, "message": "No matching handler for intent" } }`)
}

// Authorizes requests against the given audience. It also enables token verification
// and it launches a certificate retrieval goroutine. Only one proxy is accepted and if
// provided it shall be used to retrieve the certificates through it
func (this *Router) Authorize(aud string, proxy ...string) *Router {
	initialize.Do(func() {
		if len(proxy) == 1 {
			update_keys(&proxy[0])
		} else {
			update_keys(nil)
		}
	})
	this.aud = aud
	return this
}

func update_keys(proxy *string) {
	var cache_timeout int = 3600

	client := get_http_client(proxy)
	req, err := http.NewRequest("GET", google_jwks_uri, nil)

	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Couldn't retrieve certificates from google: [", err, "], token authorizations will fail")
		return
	}

	//Get the cache-control's timeout
	for _, v := range strings.Split(resp.Header.Get("cache-control"), ", ") {
		if strings.Contains(v, "max-age=") {
			age := strings.Replace(v, "max-age=", "", -1)
			seconds, err := strconv.Atoi(age)
			if err == nil {
				cache_timeout = seconds
			}
			break
		}
	}
	json.NewDecoder(resp.Body).Decode(&google_keys)

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		select {
		case <-time.After(time.Duration(cache_timeout-1) * time.Second):
			cancel()
			update_keys(proxy)
		case <-ctx.Done():
			return
		}

	}()

	log.Println("Google certificates updated; Cache's timeout is", cache_timeout)
}

func get_http_client(proxy *string) *http.Client {
	if proxy != nil {
		proxyUrl, err := url.Parse(*proxy)
		if err != nil {
			return &http.Client{}
		}

		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		return &http.Client{Transport: transport}
	}
	return &http.Client{}
}

func (this *Router) authorized(token string) bool {
	if this.aud != "" {
		token, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
			var kid string = tok.Header["kid"].(string)
			k, err := jwt.ParseRSAPublicKeyFromPEM([]byte(google_keys[kid]))
			if err != nil {
				return nil, err
			}
			return k, nil
		})

		if err != nil {
			log.Println("Received token is invalid:", err)
			return false
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["aud"] != this.aud {
				log.Println("Claims do not match this webhook", claims)
				return false
			}
		}
	}

	return true
}

// Function type for processing Google actions' fulfillment requests
// These functions must return a pointer to the response so that the Router can
// Encode it and send the proper response back.
//	func MyHandler(fr *aog.FulfillmentRequest) *aog.FulfillmentResponse {
//		return &aog.FulfillmentResponse{
//		Prompt: &aog.Prompt{
//			FirstSimple: &aog.Simple{
//				Speech: "Hello World",
//			},
//		},
//	}
type FulfillmentFunction func(fr *aog.FulfillmentRequest) *aog.FulfillmentResponse

func handleDefaultHealthCheck(fr *aog.FulfillmentRequest) *aog.FulfillmentResponse {
	return &aog.FulfillmentResponse{
		Prompt: &aog.Prompt{
			Override: true,
			FirstSimple: &aog.Simple{
				Speech: "ok",
				Text:   "",
			},
		},
	}
}

// Function to replace the default provided health check intent.
//
// It is likely this will never be required
func (this *Router) HealthCheck(fn FulfillmentFunction) {
	for _, v := range this.routes {
		if v.intent == ACTIONS_INTENT_HEALTH_CHECK {
			v.function = fn
			return
		}
	}
}

// HandleIntent creates a Route within the Router that shall be handled by fn.
//
// Additional criterias can be added to the Route by using provided or custom Matchers.
//
// Example:
//
//	myrouter.HandleIntent("hello_workld", func(fr *aog.FulfillmentRequest) *aog.FulfillmentResponse {
//		return &aog.FulfillmentResponse {
//			//...your hello world prompt goes here
//		}
//	}
func (this *Router) HandleIntent(intent string, fn FulfillmentFunction) *Route {
	rv := &Route{
		function: fn,
		intent:   intent,
		matches:  make([]Matcher, 1),
	}

	this.routes = append(this.routes, rv)

	return rv
}
