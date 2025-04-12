package main

import (
	"fmt"
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

	handler := func(w *response.Writer, req *request.Request) {

		hdrs := response.GetDefaultHeaders(0)
		hdrs.Replace("Content-Type", "text/html")

		if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(response.StatusCodeBadRequest)
			hdrs.Replace("Content-Length", fmt.Sprintf("%d", len(HTML_BAD_REQUEST)))
			w.WriteHeaders(hdrs)
			w.WriteBody([]byte(HTML_BAD_REQUEST))
			return
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			w.WriteStatusLine(response.StatusCodeInternalServerError)
			hdrs.Replace("Content-Length", fmt.Sprintf("%d", len(HTML_INTERNAL_SERVER_ERROR)))
			w.WriteHeaders(hdrs)
			w.WriteBody([]byte(HTML_INTERNAL_SERVER_ERROR))
			return
		}

		w.WriteStatusLine(response.StatusCodeOK)
		hdrs.Replace("Content-Length", fmt.Sprintf("%d", len(HTML_OK)))
		w.WriteHeaders(hdrs)
		w.WriteBody([]byte(HTML_OK))

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

const HTML_BAD_REQUEST = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const HTML_INTERNAL_SERVER_ERROR = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const HTML_OK = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
