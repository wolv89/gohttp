package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	sender, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer sender.Close()

	r := bufio.NewReader(os.Stdin)

	for {

		fmt.Print(">")
		s, err := r.ReadString('\n')
		if err != nil {
			log.Print(err)
		}

		sender.Write([]byte(s))

	}

}
