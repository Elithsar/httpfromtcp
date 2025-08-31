package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection accepted")
		go func(c net.Conn) {
			req, err := request.RequestFromReader(c)

			if err != nil {
				log.Println("Error reading request:", err)
				c.Close()
				return
			}

			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", req.RequestLine.Method)
			fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
			fmt.Println("Headers:")
			for k, v := range req.Headers {
				fmt.Printf("- %v: %v\n", k, v)
			}
			c.Close()
			fmt.Println("Connection closed")
		}(conn)
	}
}
