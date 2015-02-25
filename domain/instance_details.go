package domain

type ServiceInstanceDetailsRequest struct {
	InstanceID string
}

type ServiceInstanceDetailsResponse struct {
	DashboardURL             string
	LastOperationState       string
	LastOperationDescription string
}
