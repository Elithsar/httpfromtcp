package server

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
	"sync/atomic"
)

// HandlerError structure for consistent error handling
// Handler function type (new signature)
type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
	}
	go s.listenWithHandler(handler)
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listenWithHandler(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return // Server is closed, exit the loop
			}
			continue // Ignore other errors and continue accepting
		}
		go s.handleWithHandler(conn, handler)
	}
}

func (s *Server) handleWithHandler(conn net.Conn, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	w := response.NewWriter(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		h := headers.NewHeaders()
		h.Set("Content-Type", "text/html")
		body := []byte(`<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`)
		h.Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}

	handler(w, req)
}
