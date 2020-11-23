package main

import (
	"log"
	"protocols/transport-layer/protocol"
)

func main() {

	socket := protocol.NewSocket("udp4", ":8080")

	udpAdrr, err := socket.CreateUpdAddress()
	if err != nil {
		log.Println(err)
	}

	connection, err := socket.SocketConnect(udpAdrr)

	if err != nil {
		log.Println(err)
	}

	receivedValue := socket.Receive(connection)

	log.Println("RECEIVED VALUE:", receivedValue)

	socket.Send(receivedValue+"vova", udpAdrr)

	newVal := socket.Receive(connection)

	socket.Send(newVal, udpAdrr)

	log.Println("NEW VALUE:", newVal)

}
