package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	rAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	udpConn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		_, err = udpConn.Write([]byte(text))
	}
}
