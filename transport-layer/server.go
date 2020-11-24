package main

import (
	"log"
	"protocols/transport-layer/protocol"
)

func main() {

	log.Println("SERVER")

	socket := protocol.NewSocket("udp4", ":1234")

	udpAddress, err := socket.CreateUpdAddress()
	if err != nil {
		log.Println(err)
	}

	connection, err := socket.ServerSocketConnect(udpAddress)
	if err != nil {
		log.Println(err)
	}

	defer connection.Close()

	myValue, err := socket.ReceiveMessage(connection)
	if err != nil {
		log.Println(err)
	}

	log.Println("MY VALUES SERVER:", myValue)

	if err = socket.SendMessage("vova mtf lab", connection); err != nil {
		log.Println(err)
	}

	newVal, err := socket.ReceiveMessage(connection)

	if err != nil {
		log.Println(err)
	}

	if err = socket.SendMessage(newVal+"-no you!", connection); err != nil {
		log.Println(err)
	}

	log.Println(socket.ReceiveMessage(connection))

}
