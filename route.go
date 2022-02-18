package okgohook

import (
	"regexp"

	"github.com/ignaci0/okgohook/aog"
)

// Route stores the action for an intent along with additional conditions
// such as conditions on handlers, locale or custom provided.
//
// The when the route is matched (the intent is matched  and all the matchers
// return true, then the funcion FulfillmentFunction is executed).
type Route struct {
	function FulfillmentFunction
	intent   string
	matches  []Matcher
}

// Initializes a Route with the given function fn to the provided intent
//
//	route.New("talk_to_me",  hello_world)
func newRoute(intent string, fn FulfillmentFunction) *Route {
	return &Route{
		function: fn,
		intent:   intent,
		matches:  []Matcher{},
	}
}

func (this *Route) isMatch(r *aog.FulfillmentRequest) bool {
	if r.Intent.Name != this.intent {
		return false
	}

	if len(this.matches) == 0 {
		return true
	}

	for _, m := range this.matches {
		if m != nil && !m.Matches(r) {
			return false
		}
	}

	return true
}

// Adds a condition to match the exact handler in the aog.FulfillmentRequest
//
// Example usage:
//	myRouter.HandleIntent("hello_world", helloWolrdIntent).WithHandler("new-joiner")
func (this *Route) WithHandler(handler string) *Route {
	this.matches = append(this.matches, &handlerMatcher{exact: &handler})
	return this
}

// Adds a condition to match the provided regular expression in the aog.FulfillmentRequest
//
// Example usage:
//	myRouter.HandleIntent("hello_world", helloWolrdIntent).WithHandler(".*-joiner$")
func (this *Route) WithHandlerLike(h string) *Route {
	re := regexp.MustCompile(h)
	this.matches = append(this.matches, &handlerMatcher{re: re})
	return this
}

// Matches the locale of the fulfillment request to the provided regular expression
//
// Example usage:
//	myRouter.HandleIntent("hello_world", helloWolrdIntent).WithLocaleLike("^EN-.*")
func (this *Route) WithLocaleLike(l string) *Route {
	re := regexp.MustCompile(l)
	this.matches = append(this.matches, &localeMatcher{re: re})
	return this
}

// Adds a custom match function for the given intent to filter undesired requests
// for the given IntentHandler function
func (this *Route) MatchCustomFunc(f func(*aog.FulfillmentRequest) bool) *Route {
	this.matches = append(this.matches, &customMatcher{function: f})
	return this
}

// Takes a custom Matcher for testing the request
func (this *Route) Match(m *Matcher) *Route {
	this.matches = append(this.matches, *m)
	return this
}
