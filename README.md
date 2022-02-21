# ignaci0/okgohook

This module, composed by two packages provides the means to implement fulfillment
service webhooks so Actions On Google can send to a running server.

It's important to remark this module **does not** implement [DialogFlow fulfillments](https://developers.google.com/assistant/df-asdk/overview)
but fulfillment services for the [new conversational actions](https://developers.google.com/assistant/conversational/build).

This module's features include:

* Implementation of http.Handler interface so that it is compatible with 
http.ServeMux (and other modules such as gorilla/mux)
* It provides data types for the [fulfillment requests (complete) and response (partial)](https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill) 
; missing ones can be introduced later by comunity or myself upon need
* It implements HealthCheck system intent route
* It supports request verification by validating JWT received and the intended audience
* Additional matching rules out of the box: locale and handler
* Matcher interface to introduce custom request matchers


To get started building your own actions [this is the place to go](https://developers.google.com/assistant/conversational/build).

Since actions can be run in alfa channel, this is quite convenient to run _home server's automated tasks_.

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

When the newly created router is provided with an audience, a goroutine is launched
to retrieve and keep up to date the signing token certificates. This means the server
shall require access to the internet to retrieve them.

## TO-DOs/Roadmap

* Remove unnecessary logging and add a logger facility/middlewares
* Implement missing response types
* Change the certificates verifications to a newer keys url
* ~~Add proxy support for certificates retrieval~~
* Find a way to get rid off the aog package by autogenerating code for from the [gRPC specification](https://github.com/actions-on-google/assistant-conversation-nodejs/tree/master/src/api) 
