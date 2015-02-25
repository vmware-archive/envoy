package domain

type ProvisionRequest struct {
	InstanceID        string
	ServiceID         string
	PlanID            string
	OrganizationGUID  string
	SpaceGUID         string
	AcceptsIncomplete bool
}

type ProvisionResponse struct {
	DashboardURL             string
	Async                    bool
	LastOperationDescription string
}
