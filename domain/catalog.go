package domain

var (
	_true         = true
	_false        = false
	FreeUndefined *bool
	FreeTrue      = &_true
	FreeFalse     = &_false
)

type Catalog struct {
	Services []Service `json:"services"`
}

type Service struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Bindable        bool             `json:"bindable"`
	Plans           []Plan           `json:"plans"`
	Metadata        *ServiceMetadata `json:"metadata,omitempty"`
	Tags            []string         `json:"tags,omitempty"`
	Requires        []string         `json:"requires,omitempty"`
	DashboardClient *DashboardClient `json:"dashboard_client,omitempty"`
}

type ServiceMetadata struct {
	DisplayName         string `json:"displayName"`
	ImageURL            string `json:"imageUrl"`
	LongDescription     string `json:"longDescription"`
	ProviderDisplayName string `json:"providerDisplayName"`
	DocumentationURL    string `json:"documentationUrl"`
	SupportURL          string `json:"supportUrl"`
}

type Plan struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Free        *bool         `json:"free,omitempty"`
	Metadata    *PlanMetadata `json:"metadata,omitempty"`
}

type PlanMetadata struct {
	Bullets     []string `json:"bullets"`
	Costs       []Cost   `json:"costs"`
	DisplayName string   `json:"displayName"`
}

type Cost struct {
	Amount Amount `json:"amount"`
	Unit   string `json:"unit"`
}

type Amount struct {
	USD float64 `json:"usd"`
}

type DashboardClient struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}
