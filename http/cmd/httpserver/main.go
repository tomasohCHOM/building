package main

import (
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func myProblemYourProblem(w io.Writer, r *request.Request) *server.HandlerError {
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

func myProblemYourProblemV2(w *response.Writer, req *request.Request) {
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

func main() {
	server, err := server.Serve(port, myProblemYourProblemV2)
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
