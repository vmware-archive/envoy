package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

type binder interface {
	Bind(domain.BindRequest) (domain.BindResponse, error)
}

type BindHandler struct {
	binder
}

func NewBindHandler(binder binder) BindHandler {
	return BindHandler{
		binder: binder,
	}
}

func (handler BindHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request, err := handler.Parse(req)
	if err != nil {
		respond(w, http.StatusBadRequest, Failure{err.Error()})
		return
	}

	response, err := handler.binder.Bind(request)
	if err != nil {
		if err == domain.ServiceBindingAlreadyExistsError {
			respond(w, http.StatusConflict, EmptyJSON)
		} else {
			respond(w, http.StatusInternalServerError, Failure{
				Description: err.Error(),
			})
		}
		return
	}

	respond(w, http.StatusCreated, struct {
		Credentials    domain.BindingCredentials `json:"credentials,omitempty"`
		SyslogDrainURL string                    `json:"syslog_drain_url,omitempty"`
	}{
		Credentials:    response.Credentials,
		SyslogDrainURL: response.SyslogDrainURL,
	})
}

func (handler BindHandler) Parse(req *http.Request) (domain.BindRequest, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var params struct {
		ServiceID string `json:"service_id"`
		PlanID    string `json:"plan_id"`
		AppGUID   string `json:"app_guid"`
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return domain.BindRequest{}, errors.New("request body must be a JSON object")
	}

	expression := regexp.MustCompile(`^/v2/service_instances/(.*)/service_bindings/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)
	instanceID, bindingID := matches[1], matches[2]

	if len(instanceID) == 0 || len(bindingID) == 0 ||
		len(params.ServiceID) == 0 || len(params.PlanID) == 0 {
		return domain.BindRequest{}, errors.New("missing required field")
	}

	return domain.BindRequest{
		BindingID:  bindingID,
		InstanceID: instanceID,
		ServiceID:  params.ServiceID,
		PlanID:     params.PlanID,
		AppGUID:    params.AppGUID,
	}, nil
}
