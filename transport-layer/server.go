package main

import (
	"log"
	"protocols/transport-layer/protocol"
)

func main() {
	socket := protocol.NewSocket("udp4", ":1234")

	udpAddress, err := socket.CreateUpdAddress()
	if err != nil {
		log.Println(err)
	}

	connection, err := socket.ServerSocketConnect(udpAddress)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}

	myValue, err := socket.ReceiveMessage(connection)
	if err != nil {
		log.Println(err)
	}

	log.Println("MY VALUES SERVER:", myValue)

	if err = socket.SendMessage("vova mtf lab", udpAddress, connection); err != nil {
		log.Println(err)
	}

	newVal, err := socket.ReceiveMessage(connection)

	if err != nil {
		log.Println(err)
	}

	if err = socket.SendMessage(newVal+"-no you!", udpAddress, connection); err != nil {
		log.Println(err)
	}

	log.Println(socket.ReceiveMessage(connection))

}
