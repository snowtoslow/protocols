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

	myValue := socket.Receive(connection)

	log.Println("MY VALUES SERVER:", myValue)

	socket.Send("vova mtf lab", udpAdrr)

	newVal := socket.Receive(connection)

	socket.Send(newVal+"-no you!", udpAdrr)

	log.Println(socket.Receive(connection))

}
