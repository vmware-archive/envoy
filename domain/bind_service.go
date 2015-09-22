package domain

// BindRequest encapsulates the request payload information
// for a bind request.
type BindRequest struct {
	// BindingID is the ID value for the service binding
	// represented by this bind request.
	BindingID string

	// InstanceID is the ID value for the service instance
	// to be bound in this bind request.
	InstanceID string

	// ServiceID is the ID value of the service provided in
	// the service catalog. This service was specified when
	// the service instance was provisioned.
	ServiceID string

	// PlanID is the ID value of the plan provided in the
	// service catalog. This plan was specified when the
	// service instance was provisioned.
	PlanID string

	// AppGUID is the GUID value of the application that the
	// service instance is to be bound to in this bind request.
	AppGUID string
}

// BindResponse encapsulates the response payload information
// for a bind request.
type BindResponse struct {
	// Credentials is an open set of key-value fields used to
	// indicate credential information for this service binding.
	Credentials BindingCredentials

	// SyslogDrainURL is a URL to which CloudFoundry should
	// drain logs for the bound application.
	SyslogDrainURL string
}

// BindingCredentials is an open set of key-value fields used
// to indicate credential information for a service binding.
type BindingCredentials map[string]interface{}
