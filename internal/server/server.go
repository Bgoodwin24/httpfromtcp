package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Bgoodwin24/httpfromtcp/internal/headers"
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

func (s *Server) ProxyHandler(w *response.Writer, req *request.Request) {
	if !strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteBody([]byte("Invalid path"))
		return
	}

	targetPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	httpReq, err := http.Get(fmt.Sprintf("https://httpbin.org/%v", targetPath))
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteBody([]byte("Cannot GET from site"))
		return
	}
	defer httpReq.Body.Close()

	headers := headers.NewHeaders()
	for k, v := range httpReq.Header {
		if k != "Content-Length" && len(v) > 0 {
			headers.Set(k, v[0])
		}
	}

	headers.Set("Transfer-Encoding", "chunked")

	w.WriteStatusLine(response.StatusCode(httpReq.StatusCode))
	w.WriteHeaders(headers)

	buf := make([]byte, 1024)
	for {
		n, err := httpReq.Body.Read(buf)
		if err != nil && err != io.EOF {
			w.WriteStatusLine(response.StatusInternalServerError)
			w.WriteBody([]byte("Error reading from body"))
			return
		}
		if n == 0 {
			break
		}

		_, writeErr := w.WriteChunkedBody(buf[:n])
		if writeErr != nil {
			w.WriteStatusLine(response.StatusInternalServerError)
			w.WriteBody([]byte("Error writing chunked body"))
			return
		}
	}
	w.WriteChunkedBodyDone()
}
