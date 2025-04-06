package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wolv89/gohttp/internal/request"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {

		cn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		req, err := request.RequestFromReader(cn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		if len(req.Headers) > 0 {
			fmt.Print("\nHeaders:")
			for key, value := range req.Headers {
				fmt.Printf("\n- %s: %s", key, value)
			}
		}

		if len(req.Body) > 0 {
			fmt.Print("\nBody:\n")
			fmt.Print(string(req.Body))
		}

	}

}
