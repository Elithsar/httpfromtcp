package response

import (
	"fmt"
	"io"
)

// WriteChunkedBody writes a single chunk in chunked transfer encoding.
func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	// Write chunk size in hex + CRLF
	sizeLine := fmt.Sprintf("%x\r\n", len(p))
	n, err := io.WriteString(w.w, sizeLine)
	if err != nil {
		return n, err
	}
	// Write chunk data
	n2, err := w.w.Write(p)
	if err != nil {
		return n + n2, err
	}
	// Write trailing CRLF
	n3, err := io.WriteString(w.w, "\r\n")
	return n + n2 + n3, err
}

// WriteChunkedBodyDone writes the final zero-length chunk (trailers should follow if any).
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return io.WriteString(w.w, "0\r\n")
}
