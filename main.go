package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	tcp, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
			return
		}

		fmt.Printf("connection accepted: %v\n", conn.RemoteAddr())
		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Printf("%v\n", line)
		}
		fmt.Println("connection closed.")
	}
}

func getLinesChannel(f net.Conn) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		buf := make([]byte, 8)
		part := ""

		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					if part != "" {
						ch <- part
					}
					break
				}
				fmt.Printf("error reading connection: %v\n", err)
				break
			}

			part += string(buf[:n])
			lines := strings.Split(part, "\n")

			for i := 0; i < len(lines)-1; i++ {
				ch <- lines[i]
			}

			part = lines[len(lines)-1]
		}
	}()
	return ch
}
