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

		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			s.ProxyHandler(w, req)
			return
		}

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			htmlContent := `<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
			</html>`
			header := headers.NewHeaders()
			header.Set("Content-Type", "text/html")
			w.WriteStatusLine(response.StatusBadRequest)
			w.WriteHeaders(header)
			w.WriteBody([]byte(htmlContent))
			return
		case "/myproblem":
			htmlContent := `<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
			</html>`
			header := headers.NewHeaders()
			header.Set("Content-Type", "text/html")
			w.WriteStatusLine(response.StatusInternalServerError)
			w.WriteHeaders(header)
			w.WriteBody([]byte(htmlContent))
			return
		case "/video":
			header := headers.NewHeaders()
			header.Set("Content-Type", "video/mp4")
			data, err := os.ReadFile("./assets/vim.mp4")
			if err != nil {
				w.WriteStatusLine(response.StatusInternalServerError)
				return
			}
			w.WriteStatusLine(response.StatusOK)
			w.WriteHeaders(header)
			w.WriteBody(data)
		default:
			htmlContent := `<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
			</html>`
			header := headers.NewHeaders()
			header.Set("Content-Type", "text/html")
			w.WriteStatusLine(response.StatusOK)
			w.WriteHeaders(header)
			w.WriteBody([]byte(htmlContent))
			return
		}
	}

	go server.Serve(port, handler)

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Server gracefully stopped")
}
