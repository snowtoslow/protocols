package main

import (
	"log"
	"protocols/transport-layer/protocol"
)

func main() {

	socket := protocol.NewSocket("udp4", "127.0.0.1:1234")

	udpAddress, err := socket.CreateUpdAddress()
	if err != nil {
		log.Println(err)
	}

	connection, err := socket.ClientSocket(udpAddress)
	if err != nil {
		log.Println(err)
	}

	receivedValue, err := socket.ReceiveMessage(connection)

	if err != nil {
		log.Println(err)
	}

	defer connection.Close()

	log.Println("RECEIVED VALUE:", receivedValue)

	if err = socket.SendMessage(receivedValue+"vova", connection); err != nil {
		log.Println(err)
	}

	newVal, err := socket.ReceiveMessage(connection)

	if err != nil {
		log.Println(err)
	}

	if err = socket.SendMessage(newVal, connection); err != nil {
		log.Println(err)
	}

	log.Println("NEW VALUE:", newVal)
}
