package envoy

import "github.com/pivotal-cf-experimental/envoy/domain"

// Broker defines the interface that makes up a Service Broker for CloudFoundry.
// The Broker interface is the combined interface including all of the expected
// functionality of a service broker.
type Broker interface {
	Cataloger
	Credentialer
	Provisioner
	Binder
	Unbinder
	Deprovisioner
}

// Cataloger defines the interface for a broker component providing the catalogl
// information.
type Cataloger interface {
	Catalog() domain.Catalog
}

// Credentialer defines the interface for the Basic Auth credentials required to
// interact with the service broker.
type Credentialer interface {
	Credentials() (string, string)
}

// Provisioner defines the interface for a request to provision a service.
type Provisioner interface {
	Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error)
}

// Deprovisioner defines the interface for a request to deprovision a service.
type Deprovisioner interface {
	Deprovision(domain.DeprovisionRequest) error
}

// Binder defines the interface for a request to bind a service.
type Binder interface {
	Bind(domain.BindRequest) (domain.BindResponse, error)
}

// Unbinder defines the interface for a request to unbind a service.
type Unbinder interface {
	Unbind(domain.UnbindRequest) error
}
