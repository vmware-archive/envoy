package handlers

import (
	"net/http"
	"regexp"

	"github.com/pivotal-golang/envoy/domain"
)

type Deprovisioner interface {
	Deprovision(domain.DeprovisionRequest) error
}

type DeprovisionHandler struct {
	deprovisioner Deprovisioner
}

func NewDeprovisionHandler(deprovisioner Deprovisioner) DeprovisionHandler {
	return DeprovisionHandler{
		deprovisioner: deprovisioner,
	}
}

func (handler DeprovisionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := handler.Parse(req)

	err := handler.deprovisioner.Deprovision(request)
	if err != nil {
		switch err.(type) {
		case domain.ServiceInstanceNotFoundError:
			respond(w, http.StatusGone, EmptyJSON)
		default:
			respond(w, http.StatusInternalServerError, Failure{
				Description: err.Error(),
			})
		}
		return
	}

	respond(w, http.StatusOK, EmptyJSON)
}

func (handler DeprovisionHandler) Parse(req *http.Request) domain.DeprovisionRequest {
	expression := regexp.MustCompile(`^/v2/service_instances/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)

	return domain.DeprovisionRequest{
		InstanceID: matches[1],
		ServiceID:  req.URL.Query()["service_id"][0],
		PlanID:     req.URL.Query()["plan_id"][0],
	}
}
