package response

import (
	"bytes"
	"fmt"

	"github.com/wolv89/gohttp/internal/headers"
)

const crlf = "\r\n"

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

type Writer struct {
	response bytes.Buffer
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.response.Write(p)
}

func (w Writer) Bytes() []byte {
	return w.response.Bytes()
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if statusCode < 100 || statusCode > 600 {
		return fmt.Errorf("invalid status code: %d", statusCode)
	}

	var sl bytes.Buffer
	sl.WriteString("HTTP/1.1")

	sl.WriteString(fmt.Sprintf(" %d", statusCode))

	switch statusCode {
	case StatusCodeOK:
		sl.WriteString(" OK")
	case StatusCodeBadRequest:
		sl.WriteString(" Bad Request")
	case StatusCodeInternalServerError:
		sl.WriteString(" Internal Server Error")
	}

	sl.WriteString(crlf)

	w.Write(sl.Bytes())

	return nil

}

func GetDefaultHeaders(contentLen int) headers.Headers {

	hdrs := headers.NewHeaders()

	hdrs.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	hdrs.Set("Connection", "close")
	hdrs.Set("Content-Type", "text/plain")

	return hdrs

}

func (w *Writer) WriteHeaders(hdrs headers.Headers) error {

	if len(hdrs) == 0 {
		return fmt.Errorf("no headers provided")
	}

	for key, val := range hdrs {
		w.Write([]byte(fmt.Sprintf("%s: %s%s", key, val, crlf)))
	}

	w.Write([]byte(crlf))

	return nil

}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Write(p)
}
