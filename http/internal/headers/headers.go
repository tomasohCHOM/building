package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

var CRLF = []byte("\r\n")

var (
	ErrMalformedFieldLine = errors.New("malformed HTTP field-line")
	ErrMalformedFieldName = errors.New("malformed HTTP field name")
)

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ErrMalformedFieldLine
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasPrefix(name, []byte(" ")) || bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ErrMalformedFieldName
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	n, done := 0, false
	for {
		idx := bytes.Index(data[n:], CRLF)
		if idx == -1 {
			n = 0
			break
		}
		if idx == 0 {
			done = true
			n += len(CRLF)
			break
		}

		name, value, err := parseHeader(data[n : n+idx])
		if err != nil {
			return 0, false, err
		}
		n += idx + len(CRLF)
		h[name] = value
	}

	return n, done, nil
}
