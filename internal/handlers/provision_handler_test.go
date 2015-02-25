package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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
	Response      domain.ProvisionResponse
}

func NewProvisioner() *Provisioner {
	return &Provisioner{}
}

func (p *Provisioner) Provision(req domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	p.WasCalledWith = req

	return p.Response, p.Error
}

var _ = Describe("Provision Handler", func() {
	var handler handlers.ProvisionHandler
	var provisioner *Provisioner

	BeforeEach(func() {
		provisioner = NewProvisioner()
		handler = handlers.NewProvisionHandler(provisioner)
	})

	Context("when dashboard URL is not specified", func() {
		It("returns the correct JSON and a 201 on successful provision", func() {
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

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"last_operation": {
					"state": "succeeded"
				}
			}`))

			Expect(provisioner.WasCalledWith).To(Equal(domain.ProvisionRequest{
				InstanceID:        "i-dont-care-terribly",
				PlanID:            "my-plan-id",
				ServiceID:         "my-service-id",
				OrganizationGUID:  "my-organization-guid",
				SpaceGUID:         "my-space-guid",
				AcceptsIncomplete: false,
			}))
		})

		Context("when the last operation description is provided", func() {
			BeforeEach(func() {
				provisioner.Response.LastOperationDescription = "We provisioned all the things!"
			})

			It("returns the correct JSON and a 201 on successful provision", func() {
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

				body, err := ioutil.ReadAll(writer.Body)
				if err != nil {
					panic(err)
				}
				Expect(body).To(MatchJSON(`{
					"last_operation": {
						"state": "succeeded",
						"description": "We provisioned all the things!"
					}
				}`))
			})
		})
	})

	Context("when dashboard URL is specified", func() {
		BeforeEach(func() {
			provisioner.Response.DashboardURL = "http://www.example.com/my-silly-dashboard-url"
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

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"dashboard_url":"http://www.example.com/my-silly-dashboard-url",
				"last_operation": {
					"state": "succeeded"
				}
			}`))

			Expect(provisioner.WasCalledWith).To(Equal(domain.ProvisionRequest{
				InstanceID:        "some-other-guid",
				PlanID:            "my-plan-id",
				ServiceID:         "my-service-id",
				OrganizationGUID:  "my-organization-guid",
				SpaceGUID:         "my-space-guid",
				AcceptsIncomplete: false,
			}))
		})
	})

	Context("when the provision request is asynchronous", func() {
		It("includes the state of the accepts_incomplete field in the provision request", func() {
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

			request, err := http.NewRequest("PUT", "/v2/service_instances/some-other-guid?accepts_incomplete=true", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusCreated))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"last_operation": {
					"state": "succeeded"
				}
			}`))

			Expect(provisioner.WasCalledWith).To(Equal(domain.ProvisionRequest{
				InstanceID:        "some-other-guid",
				PlanID:            "my-plan-id",
				ServiceID:         "my-service-id",
				OrganizationGUID:  "my-organization-guid",
				SpaceGUID:         "my-space-guid",
				AcceptsIncomplete: true,
			}))
		})

		It("returns a 202 status code with correct last_operation state", func() {
			provisioner.Response.Async = true

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

			request, err := http.NewRequest("PUT", "/v2/service_instances/some-other-guid?accepts_incomplete=true", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusAccepted))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"last_operation": {
					"state": "in progress"
				}
			}`))

			Expect(provisioner.WasCalledWith).To(Equal(domain.ProvisionRequest{
				InstanceID:        "some-other-guid",
				PlanID:            "my-plan-id",
				ServiceID:         "my-service-id",
				OrganizationGUID:  "my-organization-guid",
				SpaceGUID:         "my-space-guid",
				AcceptsIncomplete: true,
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

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{"description":"BOOM!"}`))
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

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{}`))
		})
	})

	Context("when an async only provider receives a non-async request", func() {
		BeforeEach(func() {
			provisioner.Error = domain.AsyncRequiredError
		})

		It("returns a 422 and error message", func() {
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

			Expect(writer.Code).To(Equal(422))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"error": "AsyncRequired",
				"description": "This service plan requires support for asynchronous provisioning by Cloud Foundry and its clients."
			}`))
		})
	})
})
