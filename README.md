# ignaci0/okgohook

--------------------------------------------------------------------------------
This module, composed by two packages provides the means to implement webhooks
Actions On Google can send to a running server.

This module's features include:

* Implementation of http.Handler interface so that it is compatible with 
http.ServeMux (and other modules such as gorilla/mux)
* It provides data types for the [fulfillment requests (complete) and response (partial)](https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill) 
; missing ones can be introduced later by comunity or myself upon need
* It implements HealthCheck system intent route
* It supports request verification by validating JWT received and the intended audience
* Additional matching rules out of the box: locale and handler
* Matcher interface to introduce custom request matchers

---

## Install

Assuming go toolchain is available:

```sh
go get github.com/ignaci0/okgohook
```

## Examples

### Hello World

Just because:

```go
webhookRouter := okgohook.NewRouter()
webhookRouter.HandleIntent("hello", func (req *aog.FulfillmentRequest) *aog.FulfillmentResponse {
	return &aog.FulfillmentResponse {
		Prompt: &aog.Prompt {
			FirstSimple: &aog.Simple { Speech: "Hello World!" },
		},
	}
}) 

//Let's use gorilla/mux for this sample:
router := mux.NewRouter()
router.Handle("/webhook", webhookRouter)
//It plays well with other handlers, e.g.:
//router.PathPrefix("/").Handler...

srv := &http.Server {
	Handler: router,
	Address: ":8080",
	WriteTimeout: 2 * time.Second,
	ReadTimeout: 2 * time.Second,
}

log.Fatal(srv.ListenAndServe())
```

### Basic Echo

This sample fulfillment function talks back the user request:

```go
webhookRouter := okgohook.NewRouter()
webhookRouter.HandleIntent("hello", func (req *aog.FulfillmentRequest) *aog.FulfillmentResponse {
	return &aog.FulfillmentResponse {
		Prompt: &aog.Prompt {
			FirstSimple: &aog.Simple { Speech: req.Intent.Query },
		},
	}
}) 
```

### With additional matches

```go
webhookRouter := okgohook.NewRouter()
webhookRouter.HandleIntent("hello", func (req *aog.FulfillmentRequest) *aog.FulfillmentResponse { }).WithHandler("world").WithLocaleLike("EN")
webhookRouter.HandleIntent("hello", func (req *aog.FulfillmentRequest) *aog.FulfillmentResponse { }).WithHandler("world").WithLocaleLike("ES")
```

### With token verification and authorization

Currently it is not possible to verify the token without audience verification.

```go
webhookRouter := okgohook.NewRouter().Authorize("my-app")
``` 

