package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/wolv89/gohttp/internal/headers"
	"github.com/wolv89/gohttp/internal/request"
	"github.com/wolv89/gohttp/internal/response"
)

type Server struct {
	handler  Handler
	listener net.Listener
	port     int
	closed   atomic.Bool
	proxy    map[string]string
}

func Serve(port int, handler Handler, proxy map[string]string) (*Server, error) {

	server := Server{
		handler: handler,
		port:    port,
		closed:  atomic.Bool{},
		proxy:   proxy,
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

	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}

	var resp response.Writer

	if len(s.proxy) > 0 {
		for from, to := range s.proxy {
			if strings.HasPrefix(req.RequestLine.RequestTarget, from) {
				HandleProxy(conn, req, &resp, from, to)
				return
			}
		}
	}

	s.handler(&resp, req)

	conn.Write(resp.Bytes())

}

func HandleProxy(conn net.Conn, req *request.Request, resp *response.Writer, from, to string) {

	path := strings.TrimPrefix(req.RequestLine.RequestTarget, from)

	hdrs := headers.NewHeaders()
	hdrs.Set("Content-Type", "text/plain")

	pr, err := http.Get(to + path)

	if err != nil {
		resp.WriteStatusLine(response.StatusCodeInternalServerError)
		hdrs.Set("Content-Length", fmt.Sprintf("%d", len(err.Error())))
		resp.WriteHeaders(hdrs)
		resp.WriteBody([]byte(err.Error()))
		conn.Write(resp.Bytes())
		return
	}

	defer pr.Body.Close()

	resp.WriteStatusLine(response.StatusCodeOK)

	hdrs.Set("Transfer-Encoding", "chunked")

	resp.WriteHeaders(hdrs)

	for {

		buf := make([]byte, 1024)
		n, err := pr.Body.Read(buf)

		if err != nil {
			break
		}

		resp.WriteChunkedBody(buf[:n])

	}

	resp.WriteChunkedBodyDone()

	conn.Write(resp.Bytes())

}
