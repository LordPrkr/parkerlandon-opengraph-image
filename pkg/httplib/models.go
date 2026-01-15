package httplib

import "net/http"

type ResponseCode int

type Handler func(http.ResponseWriter, *http.Request) error

type Middleware func(http.Handler) http.Handler

type SuccessResponse[T any] struct {
	Ok   bool         `json:"ok"`
	Code ResponseCode `json:"code"`
	Data T            `json:"data"`
}
