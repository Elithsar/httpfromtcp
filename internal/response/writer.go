package response

import (
	"fmt"
	"io"
	"httpfromtcp/internal/headers"
)

type writerState int

const (
	stateInit writerState = iota
	stateStatus
	stateHeaders
	stateBody
)

type Writer struct {
	w           io.Writer
	state       writerState
	statusCode  StatusCode
	headers     headers.Headers
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:       w,
		state:   stateInit,
		headers: headers.NewHeaders(),
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateInit {
		return fmt.Errorf("status line must be written first")
	}
	w.statusCode = statusCode
	w.state = stateStatus
	reasonPhrase, exists := statusText[statusCode]
	if !exists {
		reasonPhrase = ""
	}
	_, err := fmt.Fprintf(w.w, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != stateStatus {
		return fmt.Errorf("headers must be written after status line")
	}
	w.state = stateHeaders
	w.headers = h
	for k, v := range h {
		if _, err := fmt.Fprintf(w.w, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w.w, "\r\n")
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateHeaders {
		return 0, fmt.Errorf("body must be written after headers")
	}
	w.state = stateBody
	return w.w.Write(p)
}
