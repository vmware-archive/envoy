package envoy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/envoy/domain"
)

func TestBrokerSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envoy Suite")
}

type TestBroker struct{}

func NewTestBroker() *TestBroker {
	return &TestBroker{}
}

func (broker TestBroker) Credentials() (string, string) {
	return "username", "password"
}

func (broker *TestBroker) Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	return domain.ProvisionResponse{}, nil
}

func (broker *TestBroker) Bind(domain.BindRequest) (domain.BindResponse, error) {
	return domain.BindResponse{}, nil
}

func (broker *TestBroker) Unbind(domain.UnbindRequest) error {
	return nil
}

func (broker *TestBroker) Deprovision(domain.DeprovisionRequest) error {
	return nil
}

func (b TestBroker) Catalog() domain.Catalog {
	return domain.Catalog{}
}

func (b TestBroker) ServiceInstanceDetails(domain.ServiceInstanceDetailsRequest) (domain.ServiceInstanceDetailsResponse, error) {
	return domain.ServiceInstanceDetailsResponse{}, nil
}
