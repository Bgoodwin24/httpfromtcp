package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	msg, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	size := make([]byte, 8)
	LINE := ""

	for {
		n, err := msg.Read(size)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("read: %s\n", LINE)
				os.Exit(0)
			}
			fmt.Printf("error reading file: %v", err)
			break
		}

		parts := strings.Split(string(size[:n]), "\n")

		for i := 0; i < len(parts)-1; i++ {
			LINE += parts[i]
			fmt.Printf("read: %s\n", LINE)
			LINE = ""
		}

		LINE += parts[len(parts)-1]
	}

}
