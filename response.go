package gouno

import "net/http"

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

// NewInternalServerErrorResponse returns a new 500 Internal Server Error response.
func NewInternalServerErrorResponse() *Response {
	return NewErrorResponse(http.StatusInternalServerError, "internal server error")
}

// NewBadRequestResponse returns a new 400 Bad Request response.
func NewBadRequestResponse() *Response {
	return NewErrorResponse(http.StatusBadRequest, "bad request")
}

// NewUnauthorizedResponse returns a new 401 Unauthorized response.
func NewUnauthorizedResponse() *Response {
	return NewErrorResponse(http.StatusUnauthorized, "unauthorized")
}

// NewForbiddenResponse returns a new 403 Forbidden response.
func NewForbiddenResponse() *Response {
	return NewErrorResponse(http.StatusForbidden, "forbidden")
}

// NewNotFoundResponse returns a new 404 Not Found response.
func NewNotFoundResponse() *Response {
	return NewErrorResponse(http.StatusNotFound, "not found")
}

// NewMethodNotAllowedResponse returns a new 405 Method Not Allowed response.
func NewMethodNotAllowedResponse() *Response {
	return NewErrorResponse(http.StatusMethodNotAllowed, "method not allowed")
}

// NewRequestTimeoutResponse returns a new 408 Request Timeout response.
func NewRequestTimeoutResponse() *Response {
	return NewErrorResponse(http.StatusRequestTimeout, "request timeout")
}

// NewConflictResponse returns a new 409 Conflict response.
func NewConflictResponse() *Response {
	return NewErrorResponse(http.StatusConflict, "conflict")
}

// NewGoneResponse returns a new 410 Gone response.
func NewGoneResponse() *Response {
	return NewErrorResponse(http.StatusGone, "gone")
}
