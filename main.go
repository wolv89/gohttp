package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {

		cn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		ch := getLinesChannel(cn)

		for line := range ch {
			fmt.Println(line)
		}

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
