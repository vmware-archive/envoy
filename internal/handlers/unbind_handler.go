package handlers

import (
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

type Unbinder interface {
	Unbind(domain.UnbindRequest) error
}

type UnbindHandler struct {
	unbinder Unbinder
}

func NewUnbindHandler(unbinder Unbinder) UnbindHandler {
	return UnbindHandler{
		unbinder: unbinder,
	}
}

func (handler UnbindHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := handler.Parse(req)

	err := handler.unbinder.Unbind(request)
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

func (handler UnbindHandler) Parse(req *http.Request) domain.UnbindRequest {
	expression := regexp.MustCompile(`^/v2/service_instances/(.*)/service_bindings/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)

	return domain.UnbindRequest{
		BindingID:  matches[2],
		InstanceID: matches[1],
		ServiceID:  req.URL.Query()["service_id"][0],
		PlanID:     req.URL.Query()["plan_id"][0],
	}
}
