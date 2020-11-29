package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"protocols/constants"
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

	buffer := make([]byte, constants.BUFF_SIZE) //4000 in case of shit

	n, add, err := connection.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
	}
	log.Println("value from server which was computed:", string(buffer[:n]), add, err)

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
