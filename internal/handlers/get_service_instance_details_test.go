package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/envoy/domain"
	"github.com/pivotal-cf-experimental/envoy/internal/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Details struct {
	WasCalledWith domain.ServiceInstanceDetailsRequest
	Response      domain.ServiceInstanceDetailsResponse
	Error         error
}

func NewDetails() *Details {
	return &Details{}
}

func (d *Details) ServiceInstanceDetails(request domain.ServiceInstanceDetailsRequest) (domain.ServiceInstanceDetailsResponse, error) {
	d.WasCalledWith = request
	return d.Response, d.Error
}

var _ = Describe("GetServiceInstanceDetails Handler", func() {
	var handler handlers.GetServiceInstanceDetailsHandler
	var details *Details

	BeforeEach(func() {
		details = NewDetails()
		handler = handlers.NewGetServiceInstanceDetailsHandler(details)
	})

	It("returns a 200 and JSON body containing the instance details", func() {
		writer := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/v2/service_instances/some-instance-guid", nil)
		if err != nil {
			panic(err)
		}
		details.Response.LastOperationState = domain.ServiceInstanceOperationSucceeded

		handler.ServeHTTP(writer, request)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
      "last_operation": {
        "state": "succeeded"
      }
    }`))

		Expect(details.WasCalledWith).To(Equal(domain.ServiceInstanceDetailsRequest{
			InstanceID: "some-instance-guid",
		}))
	})

	Context("when the dashboard URL is provided", func() {
		It("includes the dashboard URL in the response body", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/v2/service_instances/some-instance-guid", nil)
			if err != nil {
				panic(err)
			}
			details.Response.DashboardURL = "http://something.example.com"
			details.Response.LastOperationState = domain.ServiceInstanceOperationInProgress

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
        "dashboard_url": "http://something.example.com",
        "last_operation": {
          "state": "in progress"
        }
      }`))
		})
	})

	Context("when the last operation description is provided", func() {
		It("includes the description in the response body", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/v2/service_instances/some-instance-guid", nil)
			if err != nil {
				panic(err)
			}
			details.Response.LastOperationState = domain.ServiceInstanceOperationFailed
			details.Response.LastOperationDescription = "The instance is up and running!"

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.String()).To(MatchJSON(`{
        "last_operation": {
          "state": "failed",
          "description": "The instance is up and running!"
        }
      }`))
		})
	})

	Context("when the service instance does not exist", func() {
		It("returns a 410 and empty JSON body", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/v2/service_instances/some-instance-guid", nil)
			if err != nil {
				panic(err)
			}
			details.Error = domain.ServiceInstanceNotFoundError("some-instance-guid")

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusGone))
			Expect(writer.Body.String()).To(MatchJSON(`{}`))
		})
	})

	Context("when an unknown error occurs", func() {
		It("returns a 500 with error description", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/v2/service_instances/some-instance-guid", nil)
			if err != nil {
				panic(err)
			}
			details.Error = errors.New("BOOM!")

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
        "description": "BOOM!"
      }`))
		})
	})
})
