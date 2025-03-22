package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	ch := getLinesChannel(file)

	for line := range ch {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
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
				fmt.Printf("error reading file: %v\n", err)
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
