package domain

var (
	_true         = true
	_false        = false
	FreeUndefined *bool
	FreeTrue      = &_true
	FreeFalse     = &_false
)

// Catalog is the information for the services provided by a
// service broker.
type Catalog struct {
	// Services is a list of services provided by the broker.
	Services []Service `json:"services"`
}

// Service is the information for a single service provided by
// the service broker.
type Service struct {
	// ID is an identifier used to correlate this service in
	// future requests to the catalog. This must be unique within
	// CloudFoundry, using a GUID is recommended.
	ID string `json:"id"`

	// Name is the command-line-friendly name of the service that
	// will appear in the catalog. This field should be all lowercase
	// with no spaces.
	Name string `json:"name"`

	// Description is a short description of the service that will
	// appear in the catalog where displayed by CloudFoundry.
	Description string `json:"description"`

	// Bindable is used to indicate to CloudFoundry whether the service
	// can be bound to applications.
	Bindable bool `json:"bindable"`

	// Plans is a list of plans provided by this service.
	Plans []Plan `json:"plans"`

	// Metadata is a set of metadata for the service offering. This field
	// is optional.
	Metadata *ServiceMetadata `json:"metadata,omitempty"`

	// Tags is a flexible mechanism to expose a classification, attribute,
	// or base technology of a service, enabling equivalent services to be
	// swapped out without changes to dependent logic in applications,
	// buildpacks, or other services. E.g. mysql, relational, redis,
	// key-value, caching, messaging, amqp. This field is optional.
	Tags []string `json:"tags,omitempty"`

	// Requires is a list of permissions that the user would have to give
	// the service, if they provision it. This field is optional.
	Requires []string `json:"requires,omitempty"`

	// DashboardClient contains the data necessary to activate the
	// Dashboard SSO feature for this service. This field is optional.
	DashboardClient *DashboardClient `json:"dashboard_client,omitempty"`
}

// ServiceMetadata is a collection of fields that provide extra metadata
// about the service.
type ServiceMetadata struct {
	// DisplayName is the name of the service to be displayed in graphical
	// clients.
	DisplayName string `json:"displayName"`

	// ImageURL is URL to an image to be displayed along with the service.
	ImageURL string `json:"imageUrl"`

	// LongDescription is a longer description that can be displayed by
	// some clients when displaying the service information.
	LongDescription string `json:"longDescription"`

	// ProviderDisplayName is the name of the provider of the service.
	ProviderDisplayName string `json:"providerDisplayName"`

	// DocumentationURL is a URL to the documentation for the service.
	DocumentationURL string `json:"documentationUrl"`

	// SupportURL is a URL to support for the service.
	SupportURL string `json:"supportUrl"`
}

// Plan is the information for a service plan provided by the service
// broker.
type Plan struct {
	// ID is an identifier used to correlate this plan in future requests
	// to the catalog. This must be unique within Cloud Foundry, using a
	// GUID is recommended.
	ID string `json:"id"`

	// Name is the command-line-friendly name of the plan that will appear
	// in the catalog. This field should be all lowercase with no spaces.
	Name string `json:"name"`

	// Description is a short description of the service that will appear
	// in the catalog.
	Description string `json:"description"`

	// Free is the field that allows the plan to be limited by the
	// non_basic_services_allowed field in a Cloud Foundry Quota. This
	// field is optional.
	Free *bool `json:"free,omitempty"`

	// Metadata is a list of metadata for a service plan. This field is
	// optional.
	Metadata *PlanMetadata `json:"metadata,omitempty"`
}

// PlanMetadata is a collection of fields that provide extra metadata
// about the service plan.
type PlanMetadata struct {
	// Bullets are features of this plan to be displayed in a bulleted-list.
	Bullets []string `json:"bullets"`

	// Costs is list of items that describes the costs of a service, in
	// what currency, and the unit of measure.
	Costs []Cost `json:"costs"`

	// DisplayName is the name of the plan to be displayed in graphical
	// clients.
	DisplayName string `json:"displayName"`
}

// Cost is a description of the cost of a service plan.
type Cost struct {
	// Amount is a set of values describing the cost of the service plan.
	Amount Amount `json:"amount"`

	// Unit is a value describing how frequently the cost is incurred.
	Unit string `json:"unit"`
}

// Amount is a set of values describing the cost of the service plan.
type Amount map[string]float64

// DashboardClient contains the information necessary to activate the
// Dashboard SSO feature for this service.
type DashboardClient struct {
	// ID is the id of the OAuth2 client that the service intends to use.
	ID string `json:"id"`

	// Secret is a secret for the dashboard client.
	Secret string `json:"secret"`

	// RedirectURI is A domain for the service dashboard that will be
	// whitelisted by the UAA to enable SSO.
	RedirectURI string `json:"redirect_uri"`
}
