package domain

// UnbindRequest encapsulates the request payload information
// for an unbind request.
type UnbindRequest struct {
	// BindingID is the ID value for the service binding
	// represented by this unbind request.
	BindingID string

	// InstanceID is the ID value for the service instance
	// to be unbound in this unbind request.
	InstanceID string

	// ServiceID is the ID value of the service provided in
	// the service catalog. This service was specified when
	// the service instance was provisioned.
	ServiceID string

	// PlanID is the ID value of the plan provided in the
	// service catalog. This plan was specified when the
	// service instance was provisioned.
	PlanID string
}
