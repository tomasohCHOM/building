package headers

import (
	"bytes"
	"errors"
	"strings"
)

const CRLF = "\r\n"

var (
	ErrMalformedFieldLine = errors.New("malformed HTTP field-line")
	ErrMalformedFieldName = errors.New("malformed HTTP field name")
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	n, done := 0, false
	for {
		idx := bytes.Index(data[n:], []byte(CRLF))
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
		h.Set(name, value)
	}

	return n, done, nil
}

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

	if !isValidToken(string(name)) {
		return "", "", errors.New("invalid token for field name")
	}

	return string(name), string(value), nil
}

func isValidToken(token string) bool {
	if token == "" {
		return false
	}

	for _, ch := range token {
		if !(ch >= 'a' && ch <= 'z' ||
			ch >= 'A' && ch <= 'Z' ||
			ch >= '0' && ch <= '9' ||
			strings.ContainsRune("!#$%&'*+-.^_`|~", ch)) {
			return false
		}
	}
	return true
}
