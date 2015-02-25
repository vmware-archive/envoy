package handlers

import (
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

type ServiceInstanceDetailer interface {
	ServiceInstanceDetails(domain.ServiceInstanceDetailsRequest) (domain.ServiceInstanceDetailsResponse, error)
}

type ServiceInstanceDetailsHandler struct {
	detailer ServiceInstanceDetailer
}

func NewServiceInstanceDetailsHandler(detailer ServiceInstanceDetailer) ServiceInstanceDetailsHandler {
	return ServiceInstanceDetailsHandler{
		detailer: detailer,
	}
}

func (handler ServiceInstanceDetailsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := handler.Parse(req)
	response, err := handler.detailer.ServiceInstanceDetails(request)
	if err != nil {
		if _, ok := err.(domain.ServiceInstanceNotFoundError); ok {
			respond(w, http.StatusGone, EmptyJSON)
		} else {
			respond(w, http.StatusInternalServerError, Failure{
				Description: err.Error(),
			})
		}
		return
	}

	var output struct {
		DashboardURL  string `json:"dashboard_url,omitempty"`
		LastOperation struct {
			State       string `json:"state"`
			Description string `json:"description,omitempty"`
		} `json:"last_operation"`
	}

	output.DashboardURL = response.DashboardURL
	output.LastOperation.State = response.LastOperationState.String()
	output.LastOperation.Description = response.LastOperationDescription

	respond(w, http.StatusOK, output)
}

func (handler ServiceInstanceDetailsHandler) Parse(req *http.Request) domain.ServiceInstanceDetailsRequest {
	expression := regexp.MustCompile(`^/v2/service_instances/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)

	return domain.ServiceInstanceDetailsRequest{
		InstanceID: matches[1],
	}
}
