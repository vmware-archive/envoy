package domain

import "errors"

var (
	// ServiceInstanceAlreadyExistsError is an error used to
	//indicate that this service instance has already been
	// provisioned.
	ServiceInstanceAlreadyExistsError = errors.New("service instance already exists")

	// ServiceBindingAlreadyExistsError is an error used to
	// indicate that this service binding already exists.
	ServiceBindingAlreadyExistsError = errors.New("service binding already exists")
)

// ServiceInstanceNotFoundError is an error type used to indicate
// that the service instance requested for deprovisioning cannot
// be found.
type ServiceInstanceNotFoundError string

// Error returns a string representation of the error message.
func (e ServiceInstanceNotFoundError) Error() string {
	return string(e)
}

// ServiceBindingNotFoundError is an error type used to indicate
// that the service binding requested for unbinding cannot be
// found.
type ServiceBindingNotFoundError string

// Error returns a string representation of the error message.
func (s ServiceBindingNotFoundError) Error() string {
	return "The service binding was not found."
}
