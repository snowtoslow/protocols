package main

import (
	"log"
	"protocols/transport-layer/protocol"
)

func main() {

	socket := protocol.NewMagicSocket("udp4", ":1234")

	udpAddress, err := socket.CreateUdpAddress()
	if err != nil {
		log.Println("ERR0:", err)
	}

	connection, err := socket.ServerSocket(udpAddress)
	if err != nil {
		log.Println("ERR1:", err)
	}

	log.Println("SERVER!")

	defer connection.Close()

	sharedServerKey, err := socket.SendValueToClient(connection)
	if err != nil {
		log.Println(err)
	}

	log.Println("SHARED SERVER KEY!", sharedServerKey)

	for {
		/*if _,err := socket.HandleServer(connection); err != nil {
			log.Println("err", err)
		}*/

		if err := socket.SecuredReceive(connection, sharedServerKey.Bytes()); err != nil {
			log.Println(err)
		}
	}
}
