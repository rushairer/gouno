package gouno

import "net/http"

var (
	ErrInternalServerErrorResponse = NewErrorResponse(http.StatusInternalServerError, "internal server error")
	ErrBadRequestResponse          = NewErrorResponse(http.StatusBadRequest, "bad request")
	ErrUnauthorizedResponse        = NewErrorResponse(http.StatusUnauthorized, "unauthorized")
	ErrForbiddenResponse           = NewErrorResponse(http.StatusForbidden, "forbidden")
	ErrNotFoundResponse            = NewErrorResponse(http.StatusNotFound, "not found")
	ErrMethodNotAllowedResponse    = NewErrorResponse(http.StatusMethodNotAllowed, "method not allowed")
	ErrRequestTimeoutResponse      = NewErrorResponse(http.StatusRequestTimeout, "request timeout")
	ErrConflictResponse            = NewErrorResponse(http.StatusConflict, "conflict")
	ErrGoneResponse                = NewErrorResponse(http.StatusGone, "gone")
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewSuccessResponse(data interface{}) *Response {
	return NewResponse(0, "success", data)
}

func NewErrorResponse(code int, message string) *Response {
	return NewResponse(code, message, nil)
}
