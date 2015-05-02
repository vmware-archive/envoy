package handlers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/envoy/domain"
	"github.com/pivotal-cf-experimental/envoy/internal/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Cataloger struct{}

func NewCataloger() Cataloger {
	return Cataloger{}
}

func (c Cataloger) Catalog() domain.Catalog {
	return domain.Catalog{
		Services: []domain.Service{
			{
				ID:          "test-service",
				Name:        "testing",
				Description: "A testable service",
				Bindable:    true,
				Tags:        []string{"testable", "fast"},
				Metadata: &domain.ServiceMetadata{
					DisplayName:         "Testable Service",
					ImageURL:            "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg",
					LongDescription:     "This service is used to test things",
					ProviderDisplayName: "My Testing Framework",
					DocumentationURL:    "http://docs.example.com",
					SupportURL:          "http://support.example.com",
				},
				Requires: []string{"syslog_drain"},
				Plans: []domain.Plan{
					{
						ID:          "test-plan-1",
						Name:        "first",
						Description: "The first plan",
						Metadata: &domain.PlanMetadata{
							Bullets: []string{
								"Is not a real thing",
								"Runs on unicorns",
							},
							Costs: []domain.Cost{
								{
									Amount: domain.Amount{
										USD: 12.3,
									},
									Unit: "MONTHLY",
								},
							},
							DisplayName: "This plan is first",
						},
					},
					{
						ID:          "test-plan-2",
						Name:        "second",
						Description: "The second plan",
						Free:        domain.FreeFalse,
					},
				},
				DashboardClient: &domain.DashboardClient{
					ID:          "client-1",
					Secret:      "super-secret",
					RedirectURI: "http://dashboard.example.com",
				},
			},
		},
	}
}

var _ = Describe("CatalogHandler", func() {
	var handler handlers.CatalogHandler
	var cataloger Cataloger

	BeforeEach(func() {
		cataloger := NewCataloger()
		handler = handlers.NewCatalogHandler(cataloger)
	})

	It("returns a 200 status code with a JSON representation of the catalog", func() {
		writer := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/v2/catalog", nil)
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(writer, request)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Header()["Content-Type"]).To(Equal([]string{"application/json"}))

		body, err := ioutil.ReadAll(writer.Body)
		Expect(err).NotTo(HaveOccurred())

		var responseStructure domain.Catalog
		err = json.Unmarshal(body, &responseStructure)
		Expect(err).NotTo(HaveOccurred())

		Expect(responseStructure).To(Equal(cataloger.Catalog()))
	})
})
