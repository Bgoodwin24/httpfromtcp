package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(r) == 0 {
		return nil, fmt.Errorf("received empty HTTP input")
	}

	split := strings.Split(string(r), "\r\n")

	parts := strings.Fields(split[0])
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request-line, expected 3 parts but got %d", len(parts))
	}

	method := parts[0]
	request := parts[1]
	version := parts[2]

	for _, char := range method {
		if !unicode.IsUpper(char) || !unicode.IsLetter(char) {
			return nil, fmt.Errorf("selection is not the HTTP method: %s", method)
		}
	}

	if version != "HTTP/1.1" {
		return nil, fmt.Errorf("unsupported HTTP version: %s", version)
	}

	numericHTTPVersion := strings.TrimPrefix(version, "HTTP/")

	result := &Request{
		RequestLine: RequestLine{
			Method:        method,
			RequestTarget: request,
			HttpVersion:   numericHTTPVersion,
		},
	}
	return result, nil
}
