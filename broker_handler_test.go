package envoy_test

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pivotal-golang/envoy"
	"github.com/pivotal-golang/envoy/domain"
	"github.com/pivotal-golang/envoy/internal/handlers"
	"github.com/pivotal-golang/envoy/internal/middleware"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestBroker struct{}

func NewTestBroker() *TestBroker {
	return &TestBroker{}
}

func (broker TestBroker) Credentials() (string, string) {
	return "username", "password"
}

func (broker *TestBroker) Provision(instance domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	return domain.ProvisionResponse{}, nil
}

func (broker *TestBroker) Bind(binding domain.BindRequest) (domain.BindResponse, error) {
	return domain.BindResponse{}, nil
}

func (broker *TestBroker) Unbind(unbinding domain.UnbindRequest) error {
	return nil
}

func (broker *TestBroker) Deprovision(deprovision domain.DeprovisionRequest) error {
	return nil
}

func (b TestBroker) Catalog() domain.Catalog {
	return domain.Catalog{}
}

var _ = Describe("BrokerHandler", func() {
	var testBroker *TestBroker
	var router *mux.Router

	BeforeEach(func() {
		testBroker = NewTestBroker()
		router = envoy.NewBrokerHandler(testBroker).(*mux.Router)
	})

	Context("GET /v2/catalog", func() {
		It("routes to the CatalogHandler", func() {
			request, err := http.NewRequest("GET", "/v2/catalog", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeTrue())
			Expect(match.Handler).To(BeAssignableToTypeOf(middleware.Authenticator{}))
			auth := match.Handler.(middleware.Authenticator)
			Expect(auth.Handler).To(BeAssignableToTypeOf(handlers.CatalogHandler{}))
		})

		It("enforces the HTTP verb used", func() {
			request, err := http.NewRequest("POST", "/v2/catalog", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeFalse())
		})
	})

	Context("PUT /v2/service_instances/:id", func() {
		It("routes to the ProvisionHandler", func() {
			request, err := http.NewRequest("PUT", "/v2/service_instances/banana", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeTrue())
			Expect(match.Handler).To(BeAssignableToTypeOf(middleware.Authenticator{}))
			auth := match.Handler.(middleware.Authenticator)
			Expect(auth.Handler).To(BeAssignableToTypeOf(handlers.ProvisionHandler{}))
		})

		It("enforces the HTTP verb used", func() {
			request, err := http.NewRequest("GET", "/v2/service_instances/banana", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeFalse())
		})
	})

	Context("PUT /v2/service_instances/:instance_id/service_bindings/:binding_id", func() {
		It("routes to the BindHandler", func() {
			request, err := http.NewRequest("PUT", "/v2/service_instances/banana/service_bindings/panic", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeTrue())
			Expect(match.Handler).To(BeAssignableToTypeOf(middleware.Authenticator{}))
			auth := match.Handler.(middleware.Authenticator)
			Expect(auth.Handler).To(BeAssignableToTypeOf(handlers.BindHandler{}))
		})

		It("enforces the HTTP verb used", func() {
			request, err := http.NewRequest("GET", "/v2/service_instances/banana/service_bindings/panic", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeFalse())
		})
	})

	Describe("Unbind endpoint: DELETE /v2/service_instances/:instance_id/service_bindings/:binding_id", func() {
		It("routes to the UnbindHandler", func() {
			request, err := http.NewRequest("DELETE", "/v2/service_instances/my-instance/service_bindings/some-service-binding", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeTrue())
			Expect(match.Handler).To(BeAssignableToTypeOf(middleware.Authenticator{}))
			auth := match.Handler.(middleware.Authenticator)
			Expect(auth.Handler).To(BeAssignableToTypeOf(handlers.UnbindHandler{}))
		})

		It("enforces the HTTP verb used", func() {
			request, err := http.NewRequest("GET", "/v2/service_instances/my-instance/service_bindings/some-service-binding", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeFalse())
		})
	})

	Describe("Deprovision endpoint: DELETE /v2/service_instances/:instance_id", func() {
		It("routes to the DeprovisionHandler", func() {
			request, err := http.NewRequest("DELETE", "/v2/service_instances/my-instance?service_id=my-service&plan_id=my-plan", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeTrue())
			Expect(match.Handler).To(BeAssignableToTypeOf(middleware.Authenticator{}))
			auth := match.Handler.(middleware.Authenticator)
			Expect(auth.Handler).To(BeAssignableToTypeOf(handlers.DeprovisionHandler{}))
		})

		It("enforces the HTTP verb used", func() {
			request, err := http.NewRequest("GET", "/v2/service_instances/my-instance", nil)
			if err != nil {
				panic(err)
			}

			var match mux.RouteMatch
			Expect(router.Match(request, &match)).To(BeFalse())
		})
	})
})
