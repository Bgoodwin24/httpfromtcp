package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Bgoodwin24/httpfromtcp/internal/headers"
	"github.com/Bgoodwin24/httpfromtcp/internal/request"
	"github.com/Bgoodwin24/httpfromtcp/internal/response"
	"github.com/Bgoodwin24/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		s := &server.Server{}
		var statusCode response.StatusCode
		var htmlContent string

		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			s.ProxyHandler(w, req)
			return
		}

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			statusCode = response.StatusBadRequest
			htmlContent = `<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
			</html>`
		case "/myproblem":
			statusCode = response.StatusInternalServerError
			htmlContent = `<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
			</html>`
		default:
			statusCode = response.StatusOK
			htmlContent = `<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
			</html>`
		}
		headers := headers.NewHeaders()
		headers.Set("Content-Type", "text/html")

		w.WriteStatusLine(statusCode)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(htmlContent))
	}

	go server.Serve(port, handler)

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Server gracefully stopped")
}
