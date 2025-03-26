package main

import (
	"fmt"
	"net"

	"github.com/Bgoodwin24/httpfromtcp/internal/request"
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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error reading request line: %v", err)
			conn.Close()
			continue
		}

		fmt.Printf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for key, value := range req.Headers {
			fmt.Printf("- %v: %v\n", key, value)
		}
		fmt.Printf("Body:\n")
		fmt.Printf("%s", string(req.Body))
		conn.Close()
		fmt.Println("connection closed.")
	}
}
