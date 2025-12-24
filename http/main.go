package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(lines)

		curr := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if curr != "" {
					lines <- curr
					curr = ""
					break
				}
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal("oops")
			}
			buffer = buffer[:n]
			if i := bytes.IndexByte(buffer, '\n'); i != -1 {
				curr += string(buffer[:i])
				lines <- curr
				buffer = buffer[i+1:]
				curr = ""
			}
			curr += string(buffer)
		}
	}()

	return lines
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("oops")
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
