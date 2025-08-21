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

    var line string
    bytes := make([]byte, 8)
    for {
        _, err := file.Read(bytes)

	if err == io.EOF {
	    break
        } else if err != nil {
            log.Fatal(err)
        }

	parts := strings.Split(string(bytes), "\n")
	line += parts[0]
	if len(parts) > 1 {
            fmt.Println("read:", line)
	    line = parts[1]
        }
    }
}
