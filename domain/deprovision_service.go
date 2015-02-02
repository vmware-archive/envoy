package domain

type DeprovisionRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
}

type ServiceInstanceNotFoundError string

func (e ServiceInstanceNotFoundError) Error() string {
	return string(e)
}
