package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Bgoodwin24/httpfromtcp/internal/request"
	"github.com/Bgoodwin24/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := request.Parse(conn)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		conn.Close()
		return
	}

	buf := &bytes.Buffer{}

	handlerErr := s.handler(buf, req)

	if handlerErr != nil {
		err = WriteHandlerError(conn, handlerErr)
		if err != nil {
			log.Printf("Error writing handler error: %v", err)
		}
	} else {
		body := buf.Bytes()

		err = response.WriteStatusLine(conn, response.StatusOK)
		if err != nil {
			log.Printf("Error writing status line: %v", err)
			conn.Close()
			return
		}

		headers := response.GetDefaultHeaders(len(body))
		err = response.WriteHeaders(conn, headers)
		if err != nil {
			log.Printf("Error writing headers: %v", err)
			conn.Close()
			return
		}

		_, err = conn.Write(body)
		if err != nil {
			log.Printf("Error writing response body: %v", err)
			conn.Close()
			return
		}
	}
	conn.Close()
}

func WriteHandlerError(w io.Writer, herr *HandlerError) error {
	var statusCode response.StatusCode
	switch herr.StatusCode {
	case 400:
		statusCode = response.StatusBadRequest
	case 404:
		statusCode = response.StatusNotFound
	case 500:
		statusCode = response.StatusInternalServerError
	default:
		statusCode = response.StatusInternalServerError
	}

	err := response.WriteStatusLine(w, statusCode)
	if err != nil {
		return err
	}

	headers := response.GetDefaultHeaders(len(herr.Message))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, herr.Message)
	if err != nil {
		return err
	}
	return nil
}
