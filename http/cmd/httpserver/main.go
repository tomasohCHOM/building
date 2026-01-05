package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"http/internal/headers"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func handler(w io.Writer, r *request.Request) *server.HandlerError {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		return server.NewHandlerError(response.StatusBadRequest, "Your problem is not my problem\n")
	case "/myproblem":
		return server.NewHandlerError(response.StatusInternalServerError, "Woopsie, my bad\n")
	default:
		responseBody := "All good, frfr\n"
		w.Write([]byte(responseBody))
		return nil
	}
}

func handlerV2(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		responseBody := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(responseBody)

	case "/myproblem":
		w.WriteStatusLine(response.StatusInternalServerError)
		responseBody := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(responseBody)

	default:
		w.WriteStatusLine(response.StatusOK)
		responseBody := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(responseBody)
	}
}

func proxyHandler(w *response.Writer, req *request.Request) {
	if binTarget, ok := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin/"); ok {
		resp, err := http.Get("https://httpbin.org/" + binTarget)
		if err != nil {
			return
		}
		w.WriteStatusLine(response.StatusOK)
		h := response.GetDefaultHeaders(0)
		h.Remove("Content-Length")
		h.Set("Transfer-Encoding", "chunked")
		h.Set("Trailer", "X-Content-SHA256")
		h.Set("Trailer", "X-Content-Length")
		w.WriteHeaders(h)

		fullBody := []byte{}
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Fprintf(os.Stderr, "error reading response Body: %s", err)
				return
			}
			w.WriteChunkedBody(buf[:n])
			fullBody = append(fullBody, buf[:n]...)
		}

		t := headers.NewHeaders()
		t.Set("X-Content-SHA", fmt.Sprintf("%x", sha256.Sum256(fullBody)))
		t.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))

		err = w.WriteTrailers(h, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing trailers: %s", err)
		}
	}
}

func main() {
	server, err := server.Serve(port, proxyHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
