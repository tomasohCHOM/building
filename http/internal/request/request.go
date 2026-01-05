package request

import (
	"bytes"
	"errors"
	"http/internal/headers"
	"io"
	"strconv"
)

const CRLF = "\r\n"

var (
	ErrMalformedRequestLine        = errors.New("malformed http request-line")
	ErrMethodIsNotCapitalized      = errors.New("method is not capitalized")
	ErrUnsupportedHTTPVersion      = errors.New("unsupported HTTP version")
	ErrContentLengthDoesNotMatch   = errors.New("body length does not match content-length")
	ErrUnexpectedEOFBeforeComplete = errors.New("unexpected EOF before complete request")
)

type Request struct {
	RequestLine   RequestLine
	Headers       *headers.Headers
	Body          []byte
	state         parseState
	contentlength int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parseState string

const (
	StateInit    parseState = "init"
	StateHeaders parseState = "headers"
	StateBody    parseState = "body"
	StateDone    parseState = "done"
)

func NewRequest() *Request {
	return &Request{
		Headers: headers.NewHeaders(),
		Body:    []byte{},
		state:   StateInit,
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := NewRequest()
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
				if req.state != StateDone {
					return nil, ErrUnexpectedEOFBeforeComplete
				}
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
		r.state = StateHeaders

	case StateHeaders:
		n, done, err := r.Headers.Parse(data[read:])
		if err != nil {
			return 0, err
		}
		read += n
		if done {
			cl, found, err := r.getContentLength()
			if err != nil {
				return 0, err
			}
			if !found || cl == 0 {
				r.state = StateDone
				break
			}
			r.contentlength = cl
			r.state = StateBody
		}

	case StateBody:
		currentRead := data[read:]
		r.Body = append(r.Body, currentRead...)

		if len(r.Body) > r.contentlength {
			return 0, ErrContentLengthDoesNotMatch
		}
		read += len(currentRead)
		if len(r.Body) == r.contentlength {
			r.state = StateDone
		}

	case StateDone:
		break
	}
	return read, nil
}

func (r *Request) getContentLength() (int, bool, error) {
	val := r.Headers.Get("content-length")
	if val != "" {
		cl, err := strconv.Atoi(val)
		if err != nil {
			return 0, false, err
		}
		return cl, true, nil
	}
	return 0, false, nil
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

	if bytes.Equal(parts[0], bytes.ToUpper(parts[0])) {
		return nil, 0, ErrMethodIsNotCapitalized
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
	}, i + len(CRLF), nil
}
