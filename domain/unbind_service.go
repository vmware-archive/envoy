package domain

type UnbindRequest struct {
	BindingID  string
	InstanceID string
	ServiceID  string
	PlanID     string
}

type ServiceBindingNotFoundError string

func (s ServiceBindingNotFoundError) Error() string {
	return "The service binding was not found."
}
