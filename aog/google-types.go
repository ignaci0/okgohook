// The types provided in this package are for decoding and encoding
// the JSON objects set by the Actions on Google for a webhook
//
// The data types provided in this package are documented
// here: (https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill
//
// As of now, this package does not implement all the types and further work is required.
package aog

import "time"

type FulfillmentRequest struct {
	Handler Handler `json:"handler"`
	Intent  Intent  `json:"intent"`
	Scene   Scene   `json:"scene"`
	Session Session `json:"session"`
	User    User    `json:"user"`
	Home    Home    `json:"home,omitempty"`
	Device  Device  `json:"device"`
	Context Context `json:"context,omitempty"`
}

type Handler struct {
	Name string `json:"name"`
}

type Intent struct {
	Name   string                 `json:"name"`
	Params map[string]IntentParam `json:"params"`
	Query  string                 `json:"query"`
}

//https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Value
//This data type needs improvements to properly handle this
type Value interface{}

type IntentParam struct {
	Original string `json:"original"`
	Resolved string `json:"resolved"`
}

type Session struct {
	Id            string           `json:"id"`
	Params        map[string]Value `json:"params"`
	TypeOverrides []TypeOverride   `json:"typeOverrides"`
	LanguageCode  string           `json:"languageCode"`
}

type TypeOverride struct {
	Name    string      `json:"name"`
	Mode    string      `json:"mode"` //To be converted to enumerated TypeOverrideMode
	Synonym SynonymType `json:"synonym"`
}

type SynonymType struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Name     string       `json:"Name"`
	Synonyms []string     `json:"synonyms"`
	Display  EntryDisplay `json:"display"`
}

type EntryDisplay struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       Image   `json:"image"`
	Footer      string  `json:"footer"`
	OpenUrl     OpenUrl `json:"openUrl"`
}

type OpenUrl struct {
	Url  string `json:"url"`
	Hint string `json:"hint"` //To be converted to Enum of type UrlHint
}

type User struct {
	Locale               string                `json:"locale"`
	Params               map[string]Value      `json:"params,omitempty"`
	AccountLinkingStatus string                `json:"accountLinkingStatus"` //To be converted to enum AccountLinkingStatus
	VerificationStatus   string                `json:"verificationStatus"`   //To be converted to enum UserVerification status
	LastSeenTime         time.Time             `json:"lastSeenTime"`
	Engagement           Engagement            `json:"engagement"`
	PackageEntitlements  []PackageEntitlements `json:"packageEntitlements"`
	Permissions          []string              `json:"permissions"` //To be converted to enum Permission
}

type Engagement struct {
	PushNotificationIntents []IntentSubscription `json:"pushNotificationIntents"`
	DailyUpdateIntents      []IntentSubscription `json:"dailyUpdateIntents"`
}

type IntentSubscription struct {
	Intent       string `json:"intent"`
	ContentTitle string `json:"contentTitle"`
}

type PackageEntitlements struct {
	PackageName  string        `json:"packageName"`
	Entitlements []Entitlement `json:"entitlements"`
}

type Entitlement struct {
	Sku          string     `json:"sku"`
	SkuType      string     `json:"skuType"` //To be converted to enum SkuType
	InAppDetails SignedData `json:"inAppDetails"`
}

type SignedData struct {
	InAppPurchaseData  interface{} `json:"inAppPurchaseData"`
	InAppDataSignature string      `json:"inAppDataSignature"`
}

type Home struct {
	Params map[string]Value `json:"params,omitempty"`
}

type Device struct {
	Capabilities    []string `json:"capabilities"` //To convert to enum
	CurrentLocation Location `json:"currentLocation"`
	TimeZone        TimeZone `json:"timeZone"`
}

type Location struct {
	Coordinates   LatLng        `json:"coordinates"`
	PostalAddress PostalAddress `json:"postalADdress"`
}

type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PostalAddress struct {
	Revision           int      `json:"revision"`
	RegionCode         string   `json:"regionCode"`
	LanguageCode       string   `json:"languageCode"`
	PostalCode         string   `json:"postalCode"`
	SortingCode        string   `json:"sortingCode"`
	AdministrativeArea string   `json:"administrativeArea"`
	Locality           string   `json:"locality"`
	SubLocality        string   `json:"sublocality"`
	AddressLines       []string `json:"addressLines"`
	Recipients         []string `json:"recipient"`
	Organization       string   `json:"organization"`
}

type TimeZone struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type Context struct {
	Media  MediaContext  `json:"media"`
	Canvas CanvasContext `json:"canvas"`
}

type MediaContext struct {
	Progress string `json:"progress"` //formatted as ####.######s
	Index    int    `json:"index"`
}

type CanvasContext struct {
	State Value `json:"state"`
}

type Scene struct {
	Name              string          `json:"name"`
	SlotFillingStatus string          `json:"SlotFillingStatus"`
	Slots             map[string]Slot `json:"slots"`
	Next              NextScene       `json:"next"`
}

type Slot struct {
	Mode    string  `json:"mode"`   //To be converted to enum SlotMode
	Status  string  `json:"status"` //To be converted to enum SlotStatus
	Value   string  `json:"value"`  //It should be Value
	Updated bool    `json:"updated"`
	Prompt  *Prompt `json:"prompt,omitempty"`
}

type NextScene struct {
	Name string `json:"name"`
}

/*
 * Response data types
 */

type FulfillmentResponse struct {
	Prompt   *Prompt   `json:"prompt"`
	Scene    *Scene    `json:"scene,omitempty"`
	Session  *Session  `json:"session,omitempty"`
	User     *User     `json:"user,omitempty"`
	Home     *Home     `json:"home,omitempty"`
	Device   *Device   `json:"device,omitempty"`
	Expected *Expected `json:"expected,omitempty"`
}

type Prompt struct {
	Override    bool          `json:"override"`
	FirstSimple *Simple       `json:"firstSimple,omitempty"`
	Content     *Content      `json:"content,omitempty"`
	LastSimple  *Simple       `json:"lastSimple,omitempty"`
	Suggestions *[]Suggestion `json:"suggestions,omitempty"`
	Link        *Link         `json:"link,omitempty"`
	Canvas      *Canvas       `json:"canvas,omitempty"`
	OrderUpdate *OrderUpdate  `json:"orderUpdate,omitempty"`
}

type Simple struct {
	Speech string `json:"speech,omitempty"`
	Text   string `json:"text,omitempty"`
}

type Content struct {
	Media *Media `json:"media,omitempty"`
}

type Media struct {
	MediaType             string         `json:"mediaType"`
	StartOffset           string         `json:"startOffset,omitempty"`
	OptionalMediaControls []string       `json:"optionalMediaControls,omitempty"`
	MediaObjects          []*MediaObject `json:"mediaObjects"`
	RepeatMode            string         `json:"repeatMode,omitempty"` //To be a enum RepeatMode
	FirstMediaObjectIndex *int           `json:"firstMediaObjectIndex,omitempty"`
}

type MediaObject struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Url         string      `json:"url"`
	Image       *MediaImage `json:"image,omitempty"`
}

type MediaImage struct {
	Large *Image `json:"large,omitempty"`
	Icon  *Image `json:"icon,omitempty"`
}

type Image struct {
	Url    string `json:"url,omitempty"`
	Alt    string `json:"alt,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type Suggestion struct {
	Title string `json:"title"`
}

type Link struct {
	Name string  `json:"name"`
	Open OpenUrl `json:"open"`
}

type Canvas struct {
	Url                   string                `json:"url"`
	Data                  []Value               `json:"data"`
	SuppressMic           bool                  `json:"supressMic"`
	ContinuousMatchConfig ContinuousMatchConfig `json:"continuousMatchConfig"`
}

type ContinuousMatchConfig struct {
	ExpectedPhrases []ExpectedPhrase `json:"expectedPhrases"`
	DurationSeconds int              `json:"durationSeconds"`
}

type ExpectedPhrase struct {
	Phrase             string   `json:"phrase"`
	AlternativePhrases []string `json:"alternativePhrases"`
}

type OrderUpdate struct {
	//Order            Order            `json:"order"` //too large for nothing
	UpdateMask       string           `json:"updateMask"`
	UserNotification UserNotification `json:"userNotification"`
	Reason           string           `json:"reason"`
}

type UserNotification struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Expected struct {
	Speech []string `json:"speech"`
}
