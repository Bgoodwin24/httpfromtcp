package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	res, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("error resolving udp address: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, res)
	if err != nil {
		fmt.Printf("error preparing UDP connection: %v", err)
	}
	defer conn.Close()

	reader := os.Stdin

	r := bufio.NewReader(reader)

	for {
		fmt.Print("> ")
		readStr, err := r.ReadString('\n')
		if err != nil {
			log.Printf("%v", err)
		}

		w, err := conn.Write([]byte(readStr))
		if err != nil || w != len(readStr) {
			log.Printf("%v", err)
		}

	}
}
