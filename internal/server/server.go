package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/wolv89/gohttp/internal/request"
)

type Server struct {
	listener net.Listener
	port     int
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	server := Server{
		port:   port,
		closed: atomic.Bool{},
	}

	var err error

	server.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	go server.listen()

	return &server, nil

}

func (s *Server) listen() {

	for {

		cn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}

		go s.handle(cn)

	}

}

func (s *Server) Close() error {

	s.closed.Store(true)

	if s.listener != nil {
		s.listener.Close()
	}

	return nil

}

func (s *Server) handle(conn net.Conn) {

	_, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}

	resp := `HTTP/1.1 200 OK
Content-Type: text/plain

Hello World!`

	conn.Write([]byte(resp))
	conn.Close()

}
