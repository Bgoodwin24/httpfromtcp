package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	for _, c := range key {
		if !isValidHeaderChar(c) {
			return idx + 2, false, fmt.Errorf("invalid header name character: %c", c)
		}
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	loweredKey := strings.ToLower(key)
	h.Set(loweredKey, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func isValidHeaderChar(c rune) bool {
	if unicode.IsLetter(c) || unicode.IsDigit(c) {
		return true
	}

	specialChars := "!#$%&'*+-.^_`|~"
	for _, sc := range specialChars {
		if c == sc {
			return true
		}
	}

	return false
}
