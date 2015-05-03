package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

type provisioner interface {
	Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error)
}

type ProvisionHandler struct {
	provisioner
}

func NewProvisionHandler(provisioner provisioner) ProvisionHandler {
	return ProvisionHandler{
		provisioner: provisioner,
	}
}

func (handler ProvisionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request, err := handler.Parse(req)
	if err != nil {
		respond(w, http.StatusBadRequest, Failure{err.Error()})
		return
	}
	response, err := handler.provisioner.Provision(request)
	if err != nil {
		if err == domain.ServiceInstanceAlreadyExistsError {
			respond(w, http.StatusConflict, EmptyJSON)
		} else {
			respond(w, http.StatusInternalServerError, Failure{
				Description: err.Error(),
			})
		}
		return
	}

	respond(w, http.StatusCreated, struct {
		DashboardURL string `json:"dashboard_url,omitempty"`
	}{
		DashboardURL: response.DashboardURL,
	})
}

func (handler ProvisionHandler) Parse(req *http.Request) (domain.ProvisionRequest, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var params struct {
		ServiceID        string `json:"service_id"`
		PlanID           string `json:"plan_id"`
		OrganizationGUID string `json:"organization_guid"`
		SpaceGUID        string `json:"space_guid"`
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return domain.ProvisionRequest{}, errors.New("request body must be a JSON object")
	}

	expression := regexp.MustCompile(`^/v2/service_instances/(.*)$`)
	instanceID := expression.FindStringSubmatch(req.URL.Path)[1]

	if len(instanceID) == 0 || len(params.ServiceID) == 0 || len(params.PlanID) == 0 ||
		len(params.OrganizationGUID) == 0 || len(params.SpaceGUID) == 0 {
		return domain.ProvisionRequest{}, errors.New("missing required field")
	}

	return domain.ProvisionRequest{
		InstanceID:       instanceID,
		ServiceID:        params.ServiceID,
		PlanID:           params.PlanID,
		OrganizationGUID: params.OrganizationGUID,
		SpaceGUID:        params.SpaceGUID,
	}, nil
}
