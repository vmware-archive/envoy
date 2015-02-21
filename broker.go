package envoy

import "github.com/pivotal-cf-experimental/envoy/domain"

type Broker interface {
	Catalog() domain.Catalog
	Credentials() (string, string)
	Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error)
	Bind(domain.BindRequest) (domain.BindResponse, error)
	Unbind(domain.UnbindRequest) error
	Deprovision(domain.DeprovisionRequest) error
}
