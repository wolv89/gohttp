package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	ParseStatus parseStatus
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

type parseStatus int

const (
	parseStatusInitialised parseStatus = iota
	parseStatusDone
)

func RequestFromReader(reader io.Reader) (*Request, error) {

	req := Request{
		ParseStatus: parseStatusInitialised,
	}

	for req.ParseStatus != parseStatusDone {

		buf := bytes.NewBuffer(make([]byte, 0, 8))
		n, err := buf.ReadFrom(reader)

		if err != nil {
			return nil, errors.New("error reading request: " + err.Error())
		}

		p, err := req.parse(buf.Bytes())

		if err != nil {
			return nil, errors.New("error parsing request: " + err.Error())
		}

		if n == p {
			req.ParseStatus = parseStatusDone
		}

	}

	return &req, nil

}

func (r *Request) parse(data []byte) (int64, error) {

	bytesRead, requestLine, err := parseRequestLine(data)
	if err != nil {
		return 0, errors.New("error while parsing request line: " + err.Error())
	}

	if bytesRead == 0 {
		return 0, nil
	}

	r.RequestLine = requestLine
	return bytesRead, nil

}

func parseRequestLine(data []byte) (int64, RequestLine, error) {

	nextCLRF := bytes.Index(data, []byte(CRLF))
	if nextCLRF <= 0 {
		return 0, RequestLine{}, nil
	}

	line := string(data[:nextCLRF])

	parts := strings.Split(line, " ")
	requestLine := RequestLine{}

	if len(parts) != 3 {
		return 0, requestLine, errors.New("invalid request line: " + line)
	}

	requestLine.Method = parts[0]
	requestLine.RequestTarget = parts[1]
	requestLine.HttpVersion = parts[2]

	for _, m := range requestLine.Method {
		if m < 'A' || m > 'Z' {
			return 0, requestLine, errors.New("invalid request method: " + requestLine.Method)
		}
	}

	if requestLine.HttpVersion != HTTP_VERSION {
		return 0, requestLine, errors.New("invalid http version: " + requestLine.HttpVersion + ", expecting: " + HTTP_VERSION)
	}

	requestLine.HttpVersion = "1.1"

	return int64(len(line)), requestLine, nil

}
