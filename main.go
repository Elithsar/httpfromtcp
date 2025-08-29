package main

import(
    "fmt"
    "log"
    "io"
    "net"
    "strings"
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
	    ch := getLinesChannel(c)
	    for line := range(ch) {
	        fmt.Println(line)
	    }
	    c.Close()
	    fmt.Println("Connwction closed")
	}(conn)
    }
}

func getLinesChannel(f io.ReadCloser) <-chan string {
    ch := make(chan string)

    go func(){
        var line string
        bytes := make([]byte, 8, 8)
        for {
  	    _, err := f.Read(bytes)
 
 	    if err == io.EOF {
	        close(ch)
 	        break
 	    } else if err != nil {
	        log.Fatal(err)
	    }
            parts := strings.Split(string(bytes), "\n")
   	    line += parts[0]
	    if len(parts) > 1 {
	        ch <- line
	        line = parts[1]
  	    } 
        }
    }()
    return ch
}
