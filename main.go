package main

import(
    "fmt"
    "log"
    "io"
    "os"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
	log.Fatal(err)
    }

    for {
        bytes := make([]byte, 8)
        _, err := file.Read(bytes)
        
	if err == io.EOF {
	    break
        } else if err != nil {
            log.Fatal(err)
        }
	fmt.Println("read:", string(bytes))
    }

}
