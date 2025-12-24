package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("oops")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("oops")
		}
		fmt.Println("Connection has been accepted from", conn.RemoteAddr())

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("Connection has been closed!")
	}
}
