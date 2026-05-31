package gouno

import "net/http"

var (
	InternalServerErrorResponse = NewErrorResponse(http.StatusInternalServerError, "internal server error")
	BadRequestResponse          = NewErrorResponse(http.StatusBadRequest, "bad request")
	UnauthorizedResponse        = NewErrorResponse(http.StatusUnauthorized, "unauthorized")
	ForbiddenResponse           = NewErrorResponse(http.StatusForbidden, "forbidden")
	NotFoundResponse            = NewErrorResponse(http.StatusNotFound, "not found")
	MethodNotAllowedResponse    = NewErrorResponse(http.StatusMethodNotAllowed, "method not allowed")
	RequestTimeoutResponse      = NewErrorResponse(http.StatusRequestTimeout, "request timeout")
	ConflictResponse            = NewErrorResponse(http.StatusConflict, "conflict")
	GoneResponse                = NewErrorResponse(http.StatusGone, "gone")
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse(code int, message string, data any) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewSuccessResponse(data any) *Response {
	return NewResponse(0, "success", data)
}

func NewErrorResponse(code int, message string) *Response {
	return NewResponse(code, message, nil)
}
