package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrDuplicateEntry = errors.New("duplicate entry")

func IsClientError(err error) bool {
	var clientError *ClientError
	if errors.As(err, &clientError) {
		return true
	}
	return false
}

func IsServiceError(err error) bool {
	var serviceError *ServiceError
	if errors.As(err, &serviceError) {
		return true
	}
	return false
}

type ClientError struct {
	Err error
}

func NewClientError(err error) *ClientError {
	return &ClientError{
		Err: err,
	}
}

func (e *ClientError) Error() string {
	return e.Err.Error()
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

type ServiceError struct {
	Err error
}

func NewServiceError(err error) *ServiceError {
	return &ServiceError{
		Err: err,
	}
}

func (e *ServiceError) Error() string {
	return e.Err.Error()
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}
