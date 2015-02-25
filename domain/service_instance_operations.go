package domain

type ServiceInstanceOperationState int

const (
	ServiceInstanceOperationInProgress ServiceInstanceOperationState = iota
	ServiceInstanceOperationSucceeded
	ServiceInstanceOperationFailed
)

func (s ServiceInstanceOperationState) String() string {
	switch s {
	case ServiceInstanceOperationInProgress:
		return "in progress"
	case ServiceInstanceOperationSucceeded:
		return "succeeded"
	case ServiceInstanceOperationFailed:
		return "failed"
	default:
		panic("unknown state")
	}
}
