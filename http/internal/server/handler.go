package server

import (
	"http/internal/request"
	"http/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func NewHandlerError(statusCode response.StatusCode, message string) *HandlerError {
	return &HandlerError{
		StatusCode: response.StatusBadRequest,
		Message:    message,
	}
}

func (handlerErr *HandlerError) Write(w io.Writer) {
	headers := response.GetDefaultHeaders(len(handlerErr.Message))
	response.WriteStatusLine(w, handlerErr.StatusCode)
	response.WriteHeaders(w, headers)
	w.Write([]byte(handlerErr.Message))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerV2 func(w *response.Writer, req *request.Request)
