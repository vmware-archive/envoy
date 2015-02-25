package domain

import "errors"

var (
	ServiceInstanceAlreadyExistsError = errors.New("service instance already exists")
	ServiceBindingAlreadyExistsError  = errors.New("service binding already exists")
	AsyncRequiredError                = errors.New("async required")
)

type ServiceInstanceNotFoundError string

func (e ServiceInstanceNotFoundError) Error() string {
	return string(e)
}
