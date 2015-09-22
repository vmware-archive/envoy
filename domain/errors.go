package domain

var ()

// ServiceInstanceAlreadyExistsError is an error type used to
//indicate that this service instance has already been
// provisioned.
type ServiceInstanceAlreadyExistsError string

// Error returns a string representation of the error message.
func (e ServiceInstanceAlreadyExistsError) Error() string {
	return string(e)
}

// ServiceInstanceNotFoundError is an error type used to indicate
// that the service instance requested for deprovisioning cannot
// be found.
type ServiceInstanceNotFoundError string

// Error returns a string representation of the error message.
func (e ServiceInstanceNotFoundError) Error() string {
	return string(e)
}

// ServiceBindingAlreadyExistsError is an error type used to
// indicate that this service binding already exists.
type ServiceBindingAlreadyExistsError string

// Error returns a string representation of the error message.
func (e ServiceBindingAlreadyExistsError) Error() string {
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
