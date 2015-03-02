package envoy

import "github.com/pivotal-cf-experimental/envoy/domain"

type Broker interface {
	Cataloger
	Credentialer
	Provisioner
	Binder
	Unbinder
	Deprovisioner
}

type Cataloger interface {
	Catalog() domain.Catalog
}

type Credentialer interface {
	Credentials() (string, string)
}

type Provisioner interface {
	Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error)
}

type Deprovisioner interface {
	Deprovision(domain.DeprovisionRequest) error
}

type Binder interface {
	Bind(domain.BindRequest) (domain.BindResponse, error)
}

type Unbinder interface {
	Unbind(domain.UnbindRequest) error
}
