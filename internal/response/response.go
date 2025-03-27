package response

import (
	"fmt"
	"io"

	"github.com/Bgoodwin24/httpfromtcp/internal/headers"
)

type StatusCode int

type writerState int

const (
	writerStatus writerState = iota
	writerHeader
	writerBody
)

type Writer struct {
	Writer io.Writer
	state  writerState
}

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStatus {
		return fmt.Errorf("error status does not match function in use: %v", w.state)
	}
	switch statusCode {
	case StatusOK:
		fmt.Fprintf(w.Writer, "HTTP/1.1 %d OK\r\n", statusCode)
	case StatusBadRequest:
		fmt.Fprintf(w.Writer, "HTTP/1.1 %d Bad Request\r\n", statusCode)
	case StatusInternalServerError:
		fmt.Fprintf(w.Writer, "HTTP/1.1 %d Internal Server Error\r\n", statusCode)
	default:
		fmt.Fprint(w.Writer, "HTTP/1.1")
	}
	w.state = writerHeader
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.Headers{}
	header.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return header
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerHeader {
		return fmt.Errorf("error status does not match function in use: %v", w.state)
	}
	for k, v := range headers {
		_, err := fmt.Fprintf(w.Writer, "%v: %v\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.Writer, "\r\n")
	if err != nil {
		return err
	}
	w.state = writerBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerBody {
		return 0, fmt.Errorf("error status does not match function in use: %v", w.state)
	}
	return w.Writer.Write(p)
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: w,
		state:  writerStatus,
	}
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	totalBytes := 0
	contLength := fmt.Sprintf("%x\r\n", len(p))
	b, err := w.Writer.Write([]byte(contLength))
	if err != nil {
		return 0, err
	}
	totalBytes += b

	b, err = w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	totalBytes += b

	b, err = w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	totalBytes += b

	return totalBytes, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		_, err := fmt.Fprintf(w.Writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w.Writer, "\r\n")
	return err
}
