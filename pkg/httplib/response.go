package httplib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WrappedWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func WriteJSON[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

const (
	Success ResponseCode = -1 * iota
	InternalError
	ValidationError
	ParseError
	JWTError
	UnauthorizedError
	BadRequestError
	ForbiddenError
)

func NewSuccessResponse[T any](data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		Ok:   true,
		Code: Success,
		Data: data,
	}
}

type ErrorResponse struct {
	Ok      bool         `json:"ok"`
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
}

func NewErrorResponse(code ResponseCode, message string) *ErrorResponse {
	return &ErrorResponse{
		Ok:      false,
		Code:    code,
		Message: message,
	}
}

func NewInternalErrorResponse(reason string) *ErrorResponse {
	return NewErrorResponse(InternalError, fmt.Sprintf("internal server error: %s.", reason))
}

func NewValidationErrorResponse(problems map[string]string) *ErrorResponse {
	// format error string from probems map
	problemList := make([]string, len(problems))
	i := 0
	for field, problem := range problems {
		problemList[i] = fmt.Sprintf("validation error on field '%s': %s", field, problem)
		i++
	}
	err := strings.Join(problemList, "; ")

	return NewErrorResponse(ValidationError, err)
}

func NewParseErrorResponse(reason string) *ErrorResponse {
	return NewErrorResponse(ParseError, fmt.Sprintf("error while parsing: %s.", reason))
}

func NewJWTErrorResponse(err error) *ErrorResponse {
	return NewErrorResponse(JWTError, fmt.Sprintf("error while validating jwt: %s.", err.Error()))
}

func NewUnauthorizedResponse(reason string) *ErrorResponse {
	return NewErrorResponse(UnauthorizedError, fmt.Sprintf("unauthorized request: %s.", reason))
}

func NewBadRequestResponse(reason string) *ErrorResponse {
	return NewErrorResponse(BadRequestError, reason)
}

func NewForbiddenResponse(reason string) *ErrorResponse {
	return NewErrorResponse(ForbiddenError, reason)
}
