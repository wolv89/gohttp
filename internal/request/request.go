package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	HTTP_VERSION = "HTTP/1.1"
	CRLF         = "\r\n"
)

func RequestFromReader(reader io.Reader) (*Request, error) {

	rawRequest, err := io.ReadAll(reader)
	if err != nil || len(rawRequest) == 0 {
		return nil, errors.New("unable to read request: " + err.Error())
	}

	requestParts := strings.Split(string(rawRequest), CRLF)

	requestLine, err := parseRequestLine(requestParts[0])
	if err != nil {
		return nil, errors.New("unable to parse request line: " + err.Error())
	}

	req := Request{
		RequestLine: requestLine,
	}

	return &req, nil

}

func parseRequestLine(line string) (RequestLine, error) {

	parts := strings.Split(line, " ")
	requestLine := RequestLine{}

	if len(parts) != 3 {
		return requestLine, errors.New("invalid request line: " + line)
	}

	requestLine.Method = parts[0]
	requestLine.RequestTarget = parts[1]
	requestLine.HttpVersion = parts[2]

	for _, m := range requestLine.Method {
		if m < 'A' || m > 'Z' {
			return requestLine, errors.New("invalid request method: " + requestLine.Method)
		}
	}

	if requestLine.HttpVersion != HTTP_VERSION {
		return requestLine, errors.New("invalid http version: " + requestLine.HttpVersion + ", expecting: " + HTTP_VERSION)
	}

	return requestLine, nil

}
