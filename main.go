package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	messages, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Unable to open messages.txt")
	}
	defer messages.Close()

	ch := getLinesChannel(messages)

	for {
		str, status := <-ch
		if !status {
			break
		}
		fmt.Printf("read: %s\n", str)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	const READ = 8

	bytes := make([]byte, READ)
	ch := make(chan string)

	var (
		curr, line string
		err        error
		n          int
	)

	go func() {
		for {

			n, err = f.Read(bytes)
			if n == 0 || err != nil {
				break
			}

			curr = string(bytes[:n])
			parts := strings.Split(curr, "\n")

			if len(parts) <= 1 {
				line += curr
				continue
			}

			line += parts[0]

			ch <- line

			line = parts[1]

		}

		if len(line) > 0 {
			ch <- line
		}

		close(ch)
	}()

	return ch

}
