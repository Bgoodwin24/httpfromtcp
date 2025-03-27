package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bgoodwin24/httpfromtcp/internal/request"
	"github.com/Bgoodwin24/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w io.Writer, req *request.Request) *server.HandlerError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: 400,
				Message:    "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: 500,
				Message:    "Woopsie, my bad\n",
			}
		default:
			io.WriteString(w, "All good, frfr\n")
			return nil
		}
	}

	go server.Serve(port, handler)

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Server gracefully stopped")
}
