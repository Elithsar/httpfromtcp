package main

import(
    "fmt"
    "log"
    "io"
    "os"
    "strings"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
	log.Fatal(err)
    }
    defer file.Close()

    ch := getLinesChannel(file)

    for line := range(ch) {
        fmt.Println("read:", line)
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
