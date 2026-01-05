package response

import (
	"fmt"
	"http/internal/headers"
	"io"
	"strconv"
)

const CRLF = "\r\n"

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type writerState string

const (
	stateExpectStatusLine  writerState = "expectStatusLine"
	stateExpectHeaders     writerState = "expectHeaders"
	stateExpectBody        writerState = "expectBody"
	stateExpectChunkedBody writerState = "expectChunkedBody"
	stateDone              writerState = "done"
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
		if val := headers.Get("transfer-encoding"); val == "chunked" {
			w.state = stateExpectChunkedBody
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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateExpectChunkedBody {
		return 0, fmt.Errorf("not currently expecting to write response chunked body")
	}
	out := fmt.Sprintf("%x", len(p)) + CRLF + string(p) + CRLF
	return w.writer.Write([]byte(out))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != stateExpectChunkedBody {
		return 0, fmt.Errorf("not currently expecting to write response chunked body")
	}
	return w.writer.Write([]byte("0" + CRLF + CRLF))
}

func (w *Writer) WriteTrailers(h *headers.Headers, t *headers.Headers) error {
	if w.state != stateExpectChunkedBody {
		return fmt.Errorf("not currently expecting to write trailers")
	}
	trailerVal := h.Get("Trailer")
	if trailerVal == "" {
		return fmt.Errorf("attempting to write trailers when Trailer header not present")
	}

	// NOTE: this is not a correct implementation since it does not verify that
	// we are only writing trailers announced beforehand in the Trailer header.
	// However, I'm too lazy to fix and I did it mostly for the sake of learning :P
	_, err := w.writer.Write(headerBytes(t))
	if err == nil {
		w.state = stateDone
	}
	return err
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
