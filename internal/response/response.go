package response

import (
	"fmt"
	"io"

	"github.com/Bgoodwin24/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		fmt.Fprintf(w, "HTTP/1.1 %d OK\r\n", statusCode)
	case StatusBadRequest:
		fmt.Fprintf(w, "HTTP/1.1 %d Bad Request\r\n", statusCode)
	case StatusInternalServerError:
		fmt.Fprintf(w, "HTTP/1.1 %d Internal Server Error\r\n", statusCode)
	default:
		fmt.Fprint(w, "HTTP/1.1")
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.Headers{}
	header.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%v: %v\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\r\n")
	if err != nil {
		return err
	}
	return nil
}
