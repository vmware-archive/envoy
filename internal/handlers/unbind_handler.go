package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

type unbinder interface {
	Unbind(domain.UnbindRequest) error
}

type UnbindHandler struct {
	unbinder
}

func NewUnbindHandler(unbinder unbinder) UnbindHandler {
	return UnbindHandler{
		unbinder: unbinder,
	}
}

func (handler UnbindHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request, err := handler.Parse(req)
	if err != nil {
		respond(w, http.StatusBadRequest, Failure{err.Error()})
		return
	}

	err = handler.unbinder.Unbind(request)
	if err != nil {
		switch err.(type) {
		case domain.ServiceBindingNotFoundError:
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

func (handler UnbindHandler) Parse(req *http.Request) (domain.UnbindRequest, error) {
	expression := regexp.MustCompile(`^/v2/service_instances/(.*)/service_bindings/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)

	serviceIDValues := req.URL.Query()["service_id"]
	planIDValues := req.URL.Query()["plan_id"]

	if len(serviceIDValues) != 1 || len(planIDValues) != 1 {
		return domain.UnbindRequest{}, errors.New("query parameters 'service_id' and 'plan_id' are required.")
	}

	return domain.UnbindRequest{
		BindingID:  matches[2],
		InstanceID: matches[1],
		ServiceID:  serviceIDValues[0],
		PlanID:     planIDValues[0],
	}, nil
}
