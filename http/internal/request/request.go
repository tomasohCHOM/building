package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF = "\r\n"
const HTTP_VERSION = "HTTP/1.1"

func parseRequestLine(data string) (*RequestLine, error) {
	parts := strings.Split(data, " ")
	if len(parts) != 3 {
		return nil, errors.New("request line does not contain three elements")
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}

	if rl.HttpVersion != HTTP_VERSION {
		return nil, errors.New("unsupported HTTP version")
	}

	return rl, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	i := bytes.Index(data, []byte(CRLF))
	if i == -1 {
		return nil, errors.New("malformed http request-line")
	}
	rlStr := string(data[:i])
	rl, err := parseRequestLine(rlStr)
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *rl}, nil
}

