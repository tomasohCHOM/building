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

func main() {
	server, err := server.Serve(port, myProblemYourProblem)
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
