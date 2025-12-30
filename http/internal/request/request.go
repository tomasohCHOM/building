package request

import (
	"bytes"
	"errors"
	"io"
)

const CRLF = "\r\n"

var (
	ErrMalformedRequestLine   = errors.New("malformed http request-line")
	ErrUnsupportedHTTPVersion = errors.New("unsupported HTTP version")
)

type Request struct {
	RequestLine RequestLine
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parseState string

const (
	StateInit parseState = "init"
	StateDone parseState = "done"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{state: StateInit}
	buf := make([]byte, 8)
	bufLen := 0
	for req.state != StateDone {
		if bufLen == len(buf) { // resize if needed
			copyBuf := make([]byte, len(buf)*2)
			copy(copyBuf, buf)
			buf = copyBuf
		}

		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			if err == io.EOF {
				req.state = StateDone
				break
			}
			return nil, err
		}
		bufLen += n
		readN, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	switch r.state {
	case StateInit:
		rl, n, err := parseRequestLine(data[read:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		r.RequestLine = *rl
		read += n
		r.state = StateDone
	case StateDone:
		break
	}
	return read, nil
}

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
