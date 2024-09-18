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
	return fmt.Sprintf("Error %v with http status %v", err.Message, err.Status)
}
