package nop

import "github.com/pivotal-cf-experimental/envoy/domain"

// Broker provides a no-op implemenation of a service broker.
// It is meant to be used to fill in the parts of a service
// broker that may not apply to some implementations. For example,
// it could be used to provide a valid Binder and Unbinder
// implementation for a service that does not allow binding of
// services.
type Broker struct {
	Cataloger
	Credentialer
	Provisioner
	Binder
	Unbinder
	Deprovisioner
}

// Catalog provides an empty catalog.
type Cataloger struct{}

// Catalog returns the catalog representation.
func (c Cataloger) Catalog() domain.Catalog {
	return domain.Catalog{}
}

// Credentialer provides an empty set of Basic Auth credentials
// for authenticating requests to the service broker.
type Credentialer struct{}

// Credentials returns the credentials for authenticating HTTP
// requests.
func (c Credentialer) Credentials() (string, string) {
	return "", ""
}

// Provisioner provides an empty provisioning implementation.
type Provisioner struct{}

// Provision returns an empty domain.ProvisionResponse.
func (p Provisioner) Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	return domain.ProvisionResponse{}, nil
}

// Binder provides an empty binding implementation.
type Binder struct{}

// Bind returns an empty domain.BindResponse.
func (b Binder) Bind(domain.BindRequest) (domain.BindResponse, error) {
	return domain.BindResponse{}, nil
}

// Unbinder provides an empty unbinding implementation.
type Unbinder struct{}

// Unbind returns an empty domain.UnbindResponse.
func (u Unbinder) Unbind(domain.UnbindRequest) error {
	return nil
}

// Deprovisioner provides an empty deprovisioning implementation.
type Deprovisioner struct{}

// Deprovision returns an empty domain.DeprovisionResponse.
func (d Deprovisioner) Deprovision(domain.DeprovisionRequest) error {
	return nil
}
