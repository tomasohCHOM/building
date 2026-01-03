package response

import (
	"fmt"
	"http/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type writerState string

const (
	stateExpectStatusLine writerState = "expectStatusLine"
	stateExpectHeaders    writerState = "expectHeaders"
	stateExpectBody       writerState = "expectBody"
	stateDone             writerState = "done"
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer, state: stateExpectStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateExpectStatusLine {
		return fmt.Errorf("not currently expecting to write response status line")
	}
	_, err := w.writer.Write(statusLineBytes(statusCode))
	if err == nil {
		w.state = stateExpectHeaders
	}
	return err
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	if w.state != stateExpectHeaders {
		return fmt.Errorf("not currently expecting to write response headers")
	}
	_, err := w.writer.Write(headerBytes(headers))
	if err == nil {
		w.state = stateExpectBody
		if val := headers.Get("content-length"); val == "" {
			w.state = stateDone
		}
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateExpectBody {
		return 0, fmt.Errorf("not currently expecting to write response body")
	}
	n, err := w.writer.Write(p)
	if err == nil {
		w.state = stateDone
	}
	return n, err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(statusLineBytes(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	_, err := w.Write(headerBytes(headers))
	return err
}

func statusLineBytes(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		break
	}
	return fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
}

func headerBytes(headers *headers.Headers) []byte {
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	return fmt.Appendf(b, "\r\n")
}
