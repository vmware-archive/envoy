package domain

// ProvisionRequest encapsulates the request payload information
// for a provision request.
type ProvisionRequest struct {
	// InstanceID is the ID value for the service instance
	// to be provisioned in this provision request.
	InstanceID string

	// ServiceID is the ID value of the service provided in
	// the service catalog.
	ServiceID string

	// PlanID is the ID value of the plan provided in the
	// service catalog.
	PlanID string

	// OrganizationGUID is GUID value of the organization into
	// which this service instance will be provisioned.
	OrganizationGUID string

	// SpaceGUID is GUID value of the space into which this service
	// instance will be provisioned.
	SpaceGUID string
}

// ProvisionResponse encapsulates the response payload information
// for a provision request.
type ProvisionResponse struct {
	// DashboardURL is the URL of a web-based management user
	// interface for the service instance.
	DashboardURL string
}
