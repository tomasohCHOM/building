package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("error resolving address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("error setting up UDP connection: %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		read, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading string: %v\n", err)
			continue
		}
		_, err = conn.Write([]byte(read))
		if err != nil {
			fmt.Printf("error writing to connection: %v\n", err)
			continue
		}
	}
}
