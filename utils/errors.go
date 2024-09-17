package utils

import "fmt"

type HttpError struct {
	Status  int
	Message string
}

func NewHttpError(msg string, code int) *HttpError {
	return &HttpError{
		Status:  code,
		Message: msg,
	}
}

func (err *HttpError) Error() string {
	return fmt.Sprintf("%v with http status code %v", err.Message, err.Status)
}
