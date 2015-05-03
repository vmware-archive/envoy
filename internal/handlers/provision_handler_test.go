package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/envoy/domain"
	"github.com/pivotal-cf-experimental/envoy/internal/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Provisioner struct {
	WasCalledWith domain.ProvisionRequest
	Error         error
	DashboardURL  string
}

func NewProvisioner() *Provisioner {
	return &Provisioner{}
}

func (p *Provisioner) Provision(req domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	p.WasCalledWith = req

	return domain.ProvisionResponse{
		DashboardURL: p.DashboardURL,
	}, p.Error
}

var _ = Describe("Provision Handler", func() {
	var handler handlers.ProvisionHandler
	var provisioner *Provisioner

	BeforeEach(func() {
		provisioner = NewProvisioner()
		handler = handlers.NewProvisionHandler(provisioner)
	})

	Context("when dashboard URL is not specified", func() {
		It("returns empty JSON and a 201 on successful provision", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id":        "my-service-id",
				"plan_id":           "my-plan-id",
				"organization_guid": "my-organization-guid",
				"space_guid":        "my-space-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/i-dont-care-terribly", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusCreated))
			Expect(writer.Header()["Content-Type"]).To(Equal([]string{"application/json"}))

			Expect(writer.Body.String()).To(MatchJSON("{}"))

			Expect(provisioner.WasCalledWith).To(Equal(domain.ProvisionRequest{
				InstanceID:       "i-dont-care-terribly",
				PlanID:           "my-plan-id",
				ServiceID:        "my-service-id",
				OrganizationGUID: "my-organization-guid",
				SpaceGUID:        "my-space-guid",
			}))
		})
	})

	Context("when dashboard URL is specified", func() {
		BeforeEach(func() {
			provisioner.DashboardURL = "http://www.example.com/my-silly-dashboard-url"
		})

		It("returns JSON with dashboard URL and a 201 on successful provision", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id":        "my-service-id",
				"plan_id":           "my-plan-id",
				"organization_guid": "my-organization-guid",
				"space_guid":        "my-space-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/some-other-guid", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusCreated))
			Expect(writer.Header()["Content-Type"]).To(Equal([]string{"application/json"}))

			Expect(writer.Body.String()).To(MatchJSON(`{
				"dashboard_url":"http://www.example.com/my-silly-dashboard-url"
			}`))

			Expect(provisioner.WasCalledWith).To(Equal(domain.ProvisionRequest{
				InstanceID:       "some-other-guid",
				PlanID:           "my-plan-id",
				ServiceID:        "my-service-id",
				OrganizationGUID: "my-organization-guid",
				SpaceGUID:        "my-space-guid",
			}))
		})
	})

	Context("when there is a provision failure", func() {
		BeforeEach(func() {
			provisioner.Error = errors.New("BOOM!")
		})

		It("returns a 500 and the error as the body", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id":        "my-service-id",
				"plan_id":           "my-plan-id",
				"organization_guid": "my-organization-guid",
				"space_guid":        "my-space-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/some-other-guid", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Header()["Content-Type"]).To(Equal([]string{"application/json"}))

			Expect(writer.Body.String()).To(MatchJSON(`{"description":"BOOM!"}`))
		})
	})

	Context("when the service instance has already been provisioned", func() {
		BeforeEach(func() {
			provisioner.Error = domain.ServiceInstanceAlreadyExistsError
		})

		It("returns a 409 and the error message", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id":        "my-service-id",
				"plan_id":           "my-plan-id",
				"organization_guid": "my-organization-guid",
				"space_guid":        "my-space-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/a-duplicate-guid", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusConflict))
			Expect(writer.Header()["Content-Type"]).To(Equal([]string{"application/json"}))

			Expect(writer.Body.String()).To(MatchJSON(`{}`))
		})
	})
})
