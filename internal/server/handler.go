package server

import (
	"io"

	"github.com/wolv89/gohttp/internal/request"
	"github.com/wolv89/gohttp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Message    error
	StatusCode response.StatusCode
}
