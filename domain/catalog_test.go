package domain_test

import (
	"encoding/json"

	"github.com/pivotal-cf-experimental/envoy/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Catalog", func() {
	var catalog domain.Catalog

	Context("catalog with all optional fields", func() {
		var catalog domain.Catalog

		BeforeEach(func() {
			catalog = domain.Catalog{
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
		})

		It("can be correctly represented in JSON", func() {
			document, err := json.Marshal(catalog)
			Expect(err).NotTo(HaveOccurred())
			Expect(document).To(MatchJSON([]byte(`
				{
				  "services": [
					{
					  "id": "test-service",
					  "name": "testing",
					  "description": "A testable service",
					  "bindable": true,
					  "plans": [
						{
						  "id": "test-plan-1",
						  "name": "first",
						  "description": "The first plan",
						  "metadata": {
							"bullets": [
							  "Is not a real thing",
							  "Runs on unicorns"
							],
							"costs": [
							  {
								"amount": {
								  "usd": 12.3
								},
								"unit": "MONTHLY"
							  }
							],
							"displayName": "This plan is first"
						  }
						},
						{
						  "id": "test-plan-2",
						  "name": "second",
						  "description": "The second plan",
						  "free": false
						}
					  ],
					  "metadata": {
						"displayName": "Testable Service",
						"imageUrl": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg",
						"longDescription": "This service is used to test things",
						"providerDisplayName": "My Testing Framework",
						"documentationUrl": "http://docs.example.com",
						"supportUrl": "http://support.example.com"
					  },
					  "tags": [
						"testable",
						"fast"
					  ],
					  "requires": [
						"syslog_drain"
					  ],
					  "dashboard_client": {
						"id": "client-1",
						"secret": "super-secret",
						"redirect_uri": "http://dashboard.example.com"
					  }
					}
				  ]
				}

			`)))
		})
	})

	Context("catalog with only required fields", func() {
		It("can be correctly represented in JSON", func() {
			catalog = domain.Catalog{
				Services: []domain.Service{
					{
						ID:          "test-service",
						Name:        "my-test",
						Description: "Testing the catalog",
						Bindable:    false,
						Plans: []domain.Plan{
							{
								ID:          "plan-1",
								Name:        "First plan",
								Description: "this is the first plan",
							},
						},
					},
				},
			}

			document, err := json.Marshal(catalog)
			Expect(err).NotTo(HaveOccurred())
			Expect(document).To(MatchJSON([]byte(`
				{
				  "services": [
					{
					  "id": "test-service",
					  "name": "my-test",
					  "description": "Testing the catalog",
					  "bindable": false,
					  "plans": [
						{
						  "id": "plan-1",
						  "name": "First plan",
						  "description": "this is the first plan"
						}
					  ]
					}
				  ]
				}
			`)))
		})
	})
})
