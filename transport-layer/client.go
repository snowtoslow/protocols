package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"protocols/transport-layer/protocol"
	"strings"
)

func main() {

	socket := protocol.NewMagicSocket("udp", "localhost:1234")

	address, err := socket.CreateUdpAddress()
	if err != nil {
		log.Println("Err1:", err)
	}

	connection, err := socket.ClientSocket(address)

	log.Println("CLIENT!")

	if err != nil {
		log.Println("Err2", err)
	}

	fmt.Printf("The UDP server is %s\n", connection.RemoteAddr().String())
	defer connection.Close()

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(text) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}

		if err := socket.SendMessage(connection, text); err != nil {
			log.Println(err)
		}
	}

}
