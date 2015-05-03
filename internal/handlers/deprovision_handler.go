package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

type deprovisioner interface {
	Deprovision(domain.DeprovisionRequest) error
}

type DeprovisionHandler struct {
	deprovisioner
}

func NewDeprovisionHandler(deprovisioner deprovisioner) DeprovisionHandler {
	return DeprovisionHandler{
		deprovisioner: deprovisioner,
	}
}

func (handler DeprovisionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request, err := handler.Parse(req)
	if err != nil {
		respond(w, http.StatusBadRequest, Failure{err.Error()})
		return
	}

	err = handler.deprovisioner.Deprovision(request)
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

func (handler DeprovisionHandler) Parse(req *http.Request) (domain.DeprovisionRequest, error) {
	expression := regexp.MustCompile(`^/v2/service_instances/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)

	serviceIDValues := req.URL.Query()["service_id"]
	planIDValues := req.URL.Query()["plan_id"]

	if len(serviceIDValues) != 1 || len(planIDValues) != 1 {
		return domain.DeprovisionRequest{}, errors.New("query parameters 'service_id' and 'plan_id' are required.")
	}

	return domain.DeprovisionRequest{
		InstanceID: matches[1],
		ServiceID:  serviceIDValues[0],
		PlanID:     planIDValues[0],
	}, nil
}
