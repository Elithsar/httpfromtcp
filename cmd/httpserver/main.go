package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const port = 42069

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			// Proxy to httpbin.org
			path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
			url := "https://httpbin.org" + path
			resp, err := http.Get(url)
			if err != nil {
				w.WriteStatusLine(response.StatusInternalServerError)
				h := headers.NewHeaders()
				h.Set("Content-Type", "text/plain")
				w.WriteHeaders(h)
				w.WriteBody([]byte("Proxy error\n"))
				return
			}
			defer resp.Body.Close()
			w.WriteStatusLine(response.StatusCode(resp.StatusCode))
			h := headers.NewHeaders()
			for k, v := range resp.Header {
				if strings.ToLower(k) == "content-length" {
					continue
				}
				h.Set(k, v[0])
			}
			// Announce trailers
			h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
			h.Set("Transfer-Encoding", "chunked")
			w.WriteHeaders(h)
			buf := make([]byte, 1024)
			var fullBody []byte
			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					w.WriteChunkedBody(buf[:n])
					fullBody = append(fullBody, buf[:n]...)
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}
			}
			w.WriteChunkedBodyDone()
			// Compute trailers
			sum := sha256.Sum256(fullBody)
			trailerHeaders := headers.NewHeaders()
			trailerHeaders.Set("X-Content-SHA256", fmt.Sprintf("%x", sum[:]))
			trailerHeaders.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
			w.WriteTrailers(trailerHeaders)
			return
		}

		var (
			html   string
			status response.StatusCode
		)
		if req.RequestLine.RequestTarget == "/video" {
			data, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				w.WriteStatusLine(response.StatusInternalServerError)
				h := headers.NewHeaders()
				h.Set("Content-Type", "text/plain")
				w.WriteHeaders(h)
				w.WriteBody([]byte("Failed to read video file\n"))
				return
			}
			h := headers.NewHeaders()
			h.Set("Content-Type", "video/mp4")
			h.Set("Content-Length", fmt.Sprintf("%d", len(data)))
			w.WriteStatusLine(response.StatusOK)
			w.WriteHeaders(h)
			w.WriteBody(data)
			return
		}
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
