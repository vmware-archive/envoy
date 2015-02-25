package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pivotal-cf-experimental/envoy/domain"
)

const (
	StateSucceeded  = "succeeded"
	StateInProgress = "in progress"
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
		} else if err == domain.AsyncRequiredError {
			respond(w, 422, Failure{
				Description: "This service plan requires support for asynchronous provisioning by Cloud Foundry and its clients.",
				Error:       "AsyncRequired",
			})
		} else {
			respond(w, http.StatusInternalServerError, Failure{
				Description: err.Error(),
			})
		}
		return
	}
	state := StateSucceeded
	status := http.StatusCreated

	if response.Async {
		state = StateInProgress
		status = http.StatusAccepted
	}

	var output struct {
		DashboardURL  string `json:"dashboard_url,omitempty"`
		LastOperation struct {
			State       string `json:"state"`
			Description string `json:"description,omitempty"`
		} `json:"last_operation"`
	}
	output.DashboardURL = response.DashboardURL
	output.LastOperation.State = state
	output.LastOperation.Description = response.LastOperationDescription

	respond(w, status, output)
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
		InstanceID:        matches[1],
		ServiceID:         params.ServiceID,
		PlanID:            params.PlanID,
		OrganizationGUID:  params.OrganizationGUID,
		SpaceGUID:         params.SpaceGUID,
		AcceptsIncomplete: req.URL.Query().Get("accepts_incomplete") == "true",
	}
}
