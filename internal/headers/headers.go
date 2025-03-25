package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if len(data) >= 2 && data[0] == '\r' && data[1] == '\n' {
		return 2, true, nil
	}

	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return 0, false, nil
	}

	colonIndex := bytes.IndexByte(data[:crlfIndex], ':')
	if colonIndex == -1 {
		return 0, false, fmt.Errorf("invalid header format: no colon found")
	}

	keyBytes := data[:colonIndex]
	trimmedKeyBytes := bytes.TrimSpace(keyBytes)
	if len(trimmedKeyBytes) < len(keyBytes) {
		return 0, false, fmt.Errorf("invalid header format: space before colon")
	}

	valueBytes := bytes.TrimSpace(data[colonIndex+1 : crlfIndex])

	h[string(trimmedKeyBytes)] = string(valueBytes)

	return crlfIndex + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
