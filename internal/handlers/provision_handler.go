package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pivotal-golang/envoy/domain"
)

type Provisioner interface {
	Provision(domain.ProvisionRequest) (domain.ProvisionResponse, error)
}

type ProvisionHandler struct {
	provisioner Provisioner
}

func NewProvisionHandler(provisioner Provisioner) ProvisionHandler {
	return ProvisionHandler{
		provisioner: provisioner,
	}
}

func (handler ProvisionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := handler.Parse(req)
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

func (handler ProvisionHandler) Parse(req *http.Request) domain.ProvisionRequest {
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
		panic(err)
	}

	expression := regexp.MustCompile(`^/v2/service_instances/(.*)$`)
	matches := expression.FindStringSubmatch(req.URL.Path)

	return domain.ProvisionRequest{
		InstanceID:       matches[1],
		ServiceID:        params.ServiceID,
		PlanID:           params.PlanID,
		OrganizationGUID: params.OrganizationGUID,
		SpaceGUID:        params.SpaceGUID,
	}
}