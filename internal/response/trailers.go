package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
)

// WriteTrailers writes HTTP chunked trailers (headers at the end of chunked body)
func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		if _, err := fmt.Fprintf(w.w, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w.w, "\r\n")
	return err
}
