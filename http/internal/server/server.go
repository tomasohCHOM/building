package server

import (
	"bytes"
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener, handler: handler}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := NewHandlerError(response.StatusBadRequest, err.Error())
		handlerError.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
}
