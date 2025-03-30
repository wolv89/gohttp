package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	messages, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Unable to open messages.txt")
	}
	defer messages.Close()

	const READ = 8

	bytes := make([]byte, READ)
	var n int

	for {

		n, err = messages.Read(bytes)

		if n < READ {
			bytes = make([]byte, n)
			messages.Seek(int64(-n), io.SeekCurrent)
			n, err = messages.Read(bytes)
		}

		if n == 0 || err != nil {
			break
		}

		fmt.Printf("read: %s\n", bytes)

	}

}
