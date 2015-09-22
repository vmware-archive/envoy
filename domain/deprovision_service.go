package domain

// DeprovisionRequest encapsulates the request payload for a
// deprovision request.
type DeprovisionRequest struct {
	// InstanceID is the ID value for the service instance to
	// be deprovisioned.
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
