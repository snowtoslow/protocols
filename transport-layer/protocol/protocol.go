package protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"protocols/constants"
	"protocols/models"
	"protocols/utils"
)

type Socket struct {
	networkType string
	port        string
}

func NewSocket(networkType string, port string) *Socket {
	return &Socket{
		networkType: networkType,
		port:        port,
	}
}

// need to split
func (socket *Socket) Send(message string, udpAddress *net.UDPAddr) {
	buffer := make([]byte, constants.BUFF_SIZE)
	var receivedStruct *models.Packet

	conn, err := socket.SocketConnect(udpAddress)
	if err != nil {
		log.Println(err)
	}

	myMagicBytes := fmt.Sprintf("%v", utils.CreatePacket(message))

	_, err = conn.WriteToUDP([]byte(myMagicBytes), udpAddress)
	if err != nil {
		log.Println(err)
	}
	_, address, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(buffer, &receivedStruct)
	if err != nil {
		log.Println(err)
	}

	log.Println("MY RECEIVED STRUCT!", receivedStruct)

	if receivedStruct.Payload != "nack" {
		_, err := conn.WriteToUDP([]byte(myMagicBytes), address)
		if err != nil {
			log.Println(err)
		}
		_, _, err = conn.ReadFromUDP(buffer)

		if err != nil {
			log.Println(err)
		}
	}

}

func (socket *Socket) Receive(conn *net.UDPConn) (valToRet string) {
	var receivedStruct *models.Packet

	buffer := make([]byte, constants.BUFF_SIZE)
	for {

		_, address, err := conn.ReadFromUDP(buffer)

		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(buffer, &receivedStruct)
		if err != nil {
			log.Println("1.", err)
		}

		if utils.ValidatePacket(receivedStruct) {
			_, err := conn.WriteToUDP([]byte("ack"), address)
			if err != nil {
				log.Println("2.", err)
			}

			if receivedStruct.Payload == "connect" {
				//socket.Port = address
				return ""
			} else {
				valToRet = receivedStruct.Payload
			}
		} else {
			myMagicBytes := fmt.Sprintf("%v", utils.CreatePacket("nack"))

			_, err = conn.WriteToUDP([]byte(myMagicBytes), address)
		}
		return
	}

}

// Listen at selected port!;
func (socket *Socket) SocketConnect(udpAddress *net.UDPAddr) (conn *net.UDPConn, err error) {

	conn, err = net.ListenUDP(socket.networkType, udpAddress)

	if err != nil {
		return nil, err
	}

	return conn, err
}

//Create a pointer to udp address;
func (socket *Socket) CreateUpdAddress() (udpAddress *net.UDPAddr, err error) {
	udpAddress, err = net.ResolveUDPAddr(socket.networkType, socket.port)
	if err != nil {
		return nil, err
	}

	return udpAddress, nil
}
