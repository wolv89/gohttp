package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wolv89/gohttp/internal/request"
	"github.com/wolv89/gohttp/internal/response"
	"github.com/wolv89/gohttp/internal/server"
)

const port = 42069

func main() {

	handler := func(w io.Writer, req *request.Request) *server.HandlerError {

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				Message:    fmt.Errorf("Your problem is not my problem\n"),
				StatusCode: response.StatusCodeBadRequest,
			}
		case "/myproblem":
			return &server.HandlerError{
				Message:    fmt.Errorf("Woopsie, my bad\n"),
				StatusCode: response.StatusCodeInternalServerError,
			}
		}

		w.Write([]byte("All good, frfr\n"))

		return nil

	}

	server, err := server.Serve(port, handler)
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
