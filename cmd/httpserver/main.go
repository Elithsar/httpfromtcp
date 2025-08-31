package main

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
)

const port = 42069

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		var (
			html   string
			status response.StatusCode
		)
		h := headers.NewHeaders()
		h.Set("Content-Type", "text/html")
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			status = response.StatusBadRequest
			html = `<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`
		case "/myproblem":
			status = response.StatusInternalServerError
			html = `<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>`
		default:
			status = response.StatusOK
			html = `<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>`
		}
		h.Set("Content-Length", fmt.Sprintf("%d", len(html)))
		w.WriteStatusLine(status)
		w.WriteHeaders(h)
		w.WriteBody([]byte(html))
	}

	srv, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Close()

	select {} // Block forever
}
