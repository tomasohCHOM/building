package server

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func NewHandlerError(message string) *HandlerError {
	return &HandlerError{
		StatusCode: response.StatusBadRequest,
		Message:    message,
	}
}

func (handlerErr *HandlerError) Write(w io.Writer) {
	w.Write([]byte(fmt.Sprintf("status %s: %s", handlerErr.StatusCode, handlerErr.Message)))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
