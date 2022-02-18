// Provides the data types for the routes matchinng
package okgohook

import (
	"regexp"

	"github.com/ignaci0/okgohook/aog"
)

type customMatcher struct {
	function func(*aog.FulfillmentRequest) bool
}

func (this *customMatcher) Matches(fr *aog.FulfillmentRequest) bool {
	return this.function(fr)
}

// A okgohook.Matcher can be used as a custom match for routes to match
// a given request. If the request is not to be processed by the route
// Matches shall return false, otherwise (yes or condition ignored) true.
type Matcher interface {
	Matches(*aog.FulfillmentRequest) bool
}

type handlerMatcher struct {
	exact *string
	re    *regexp.Regexp
}

func (this *handlerMatcher) Matches(fr *aog.FulfillmentRequest) bool {
	if fr == nil {
		//Nothing to execute
		return false
	}

	if this.exact != nil {
		return fr.Handler.Name == *this.exact
	}

	if this.re != nil {
		return this.re.FindString(fr.Handler.Name) != ""
	}

	return false
}

type localeMatcher struct {
	re *regexp.Regexp
}

func (this *localeMatcher) Matches(fr *aog.FulfillmentRequest) bool {
	if fr == nil {
		//Nothing to execute
		return false
	}

	if this.re != nil {
		return this.re.FindString(fr.User.Locale) != ""
	}

	return false
}
