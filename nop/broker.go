package nop

import "github.com/pivotal-cf-experimental/envoy/domain"

type Broker struct {
	Cataloger
	Credentialer
	Provisioner
	Binder
	Unbinder
	Deprovisioner
	ServiceInstanceDetailer
}

type Cataloger struct{}

func (c Cataloger) Catalog() domain.Catalog {
	return domain.Catalog{}
}

type Credentialer struct{}

func (c Credentialer) Credentials() (string, string) {
	return "", ""
}

type Provisioner struct{}

func (p Provisioner) Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	return domain.ProvisionResponse{}, nil
}

type Binder struct{}

func (b Binder) Bind(domain.BindRequest) (domain.BindResponse, error) {
	return domain.BindResponse{}, nil
}

type Unbinder struct{}

func (u Unbinder) Unbind(domain.UnbindRequest) error {
	return nil
}

type Deprovisioner struct{}

func (d Deprovisioner) Deprovision(domain.DeprovisionRequest) error {
	return nil
}

type ServiceInstanceDetailer struct{}

func (sid ServiceInstanceDetailer) ServiceInstanceDetails(domain.ServiceInstanceDetailsRequest) (domain.ServiceInstanceDetailsResponse, error) {
	return domain.ServiceInstanceDetailsResponse{}, nil
}
