package errors

import (
	"net/http"
)

type RestErr struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"code,omitempty"`
	Error      string `json:"error,omitempty"`
}

func NewBadRequest(message string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Error:      "bad_request",
	}
}

func NewNotFound(message string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Error:      "not_found",
	}
}
