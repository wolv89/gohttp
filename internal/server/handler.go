package server

import (
	"github.com/wolv89/gohttp/internal/request"
	"github.com/wolv89/gohttp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
