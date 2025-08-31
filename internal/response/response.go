package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
	ErrInvalidStatusCode                 = "invalid status code"
)

var statusText = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase, exists := statusText[statusCode]
	if !exists {
		reasonPhrase = ""
	}
	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}
	// End of headers
	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}
	return nil
}

func ValidateStatusCode(code int) (StatusCode, error) {
	switch StatusCode(code) {
	case StatusOK, StatusBadRequest, StatusInternalServerError:
		return StatusCode(code), nil
	default:
		return 0, errors.New(ErrInvalidStatusCode)
	}
}
