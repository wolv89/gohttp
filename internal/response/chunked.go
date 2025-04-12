package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	return w.Write([]byte(fmt.Sprintf("%x\r\n%s\r\n", len(p), p)))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Write([]byte("0\r\n"))
}
