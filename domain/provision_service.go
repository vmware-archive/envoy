package domain

type ProvisionRequest struct {
	InstanceID       string
	ServiceID        string
	PlanID           string
	OrganizationGUID string
	SpaceGUID        string
}

type ProvisionResponse struct {
	DashboardURL string
}
