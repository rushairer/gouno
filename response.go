package gouno

import "net/http"

var (
	// InternalServerErrorResponse is a preset 500 Internal Server Error response.
	InternalServerErrorResponse = NewErrorResponse(http.StatusInternalServerError, "internal server error")
	// BadRequestResponse is a preset 400 Bad Request response.
	BadRequestResponse = NewErrorResponse(http.StatusBadRequest, "bad request")
	// UnauthorizedResponse is a preset 401 Unauthorized response.
	UnauthorizedResponse = NewErrorResponse(http.StatusUnauthorized, "unauthorized")
	// ForbiddenResponse is a preset 403 Forbidden response.
	ForbiddenResponse = NewErrorResponse(http.StatusForbidden, "forbidden")
	// NotFoundResponse is a preset 404 Not Found response.
	NotFoundResponse = NewErrorResponse(http.StatusNotFound, "not found")
	// MethodNotAllowedResponse is a preset 405 Method Not Allowed response.
	MethodNotAllowedResponse = NewErrorResponse(http.StatusMethodNotAllowed, "method not allowed")
	// RequestTimeoutResponse is a preset 408 Request Timeout response.
	RequestTimeoutResponse = NewErrorResponse(http.StatusRequestTimeout, "request timeout")
	// ConflictResponse is a preset 409 Conflict response.
	ConflictResponse = NewErrorResponse(http.StatusConflict, "conflict")
	// GoneResponse is a preset 410 Gone response.
	GoneResponse = NewErrorResponse(http.StatusGone, "gone")
)

// Response represents a unified JSON API response with a status code, message, and optional data.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewResponse creates a new Response with the given HTTP status code, message, and optional data.
func NewResponse(code int, message string, data any) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponse creates a 200 OK response with the given data and message "success".
func NewSuccessResponse(data any) *Response {
	return NewResponse(http.StatusOK, "success", data)
}

// NewErrorResponse creates an error response with the given HTTP status code and message.
func NewErrorResponse(code int, message string) *Response {
	return NewResponse(code, message, nil)
}
