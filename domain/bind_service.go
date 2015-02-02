package domain

type BindRequest struct {
	BindingID  string
	InstanceID string
	ServiceID  string
	PlanID     string
	AppGUID    string
}

type BindResponse struct {
	Credentials    BindingCredentials
	SyslogDrainURL string
}

type BindingCredentials map[string]interface{}
