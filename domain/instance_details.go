package domain

const (
	ServiceInstanceOperationSucceeded  = "succeeded"
	ServiceInstanceOperationFailed     = "failed"
	ServiceInstanceOperationInProgress = "in progress"
)

type ServiceInstanceDetailsRequest struct {
	InstanceID string
}

type ServiceInstanceDetailsResponse struct {
	DashboardURL             string
	LastOperationState       string
	LastOperationDescription string
}
