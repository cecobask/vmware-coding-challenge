package entity

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type ResponseError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Details    string `json:"details"`
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf("error message = %s; status code = %d; details = %s", err.Message, err.StatusCode, err.Details)
}

func (err *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, err.StatusCode)
	return nil
}

func NewResponseError(message string, statusCode int, details string) *ResponseError {
	return &ResponseError{
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

func ErrorBadRequest(details string) *ResponseError {
	return NewResponseError("bad request error", http.StatusBadRequest, details)
}

func ErrorInternalServer(details string) *ResponseError {
	return NewResponseError("internal server error", http.StatusInternalServerError, details)
}
