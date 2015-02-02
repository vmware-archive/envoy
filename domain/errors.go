package domain

import "errors"

var (
	ServiceInstanceAlreadyExistsError = errors.New("service instance already exists")
	ServiceBindingAlreadyExistsError  = errors.New("service binding already exists")
)
