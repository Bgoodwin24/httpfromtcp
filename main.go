package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	msg, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	size := make([]byte, 8)

	for {
		n, err := msg.Read(size)
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			fmt.Printf("error reading file: %v", err)
			break
		}

		fmt.Printf("read: %s\n", string(size[:n]))
	}

}
