package server

import (
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

type Handler func(w *response.Writer, req *request.Request)

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

	writer := &response.Writer{
		Writer: conn,
	}
	s.handler(writer, req)

	conn.Close()
}
