package request

import (
	"bytes"
	"errors"
	"io"
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

var (
	ErrMalformedRequestLine   = errors.New("malformed http request-line")
	ErrUnsupportedHTTPVersion = errors.New("unsupported HTTP version")
)

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	i := bytes.Index(data, []byte(CRLF))
	if i == -1 {
		return nil, 0, nil
	}
	parts := bytes.Split(data[:i], []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	httpVersionParts := bytes.Split(parts[2], []byte("/"))
	if len(httpVersionParts) != 2 ||
		string(httpVersionParts[0]) != "HTTP" ||
		string(httpVersionParts[1]) != "1.1" {
		return nil, 0, ErrUnsupportedHTTPVersion
	}

	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpVersionParts[1]),
	}, i, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	rl, _, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *rl}, nil
}
