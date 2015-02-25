package domain

type DeprovisionRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
}
