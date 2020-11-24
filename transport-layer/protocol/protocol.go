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
func (socket *Socket) ReceiveMessage(connection *net.UDPConn) (valToRet string, err error) {
	var receivedStruct *models.Packet
	var myMagicBytes string

	buffer := make([]byte, constants.BUFF_SIZE)
	for {

		_, _, err := connection.ReadFromUDP(buffer)

		if err != nil {
			return "", err
		}

		err = json.Unmarshal(buffer, &receivedStruct)
		if err != nil {
			return "", err
		}

		if utils.ValidatePacket(receivedStruct) {
			myMagicBytes = fmt.Sprintf("%v", utils.CreatePacket("ack"))
			_, err = connection.Write([]byte(myMagicBytes))
			if err != nil {
				return "", err
			}

			if receivedStruct.Payload == "connect" {
				//socket.port = address.Port
				log.Println("connection established!")
				return "connection established!", nil
			} else {
				valToRet = receivedStruct.Payload
			}
		} else {
			myMagicBytes = fmt.Sprintf("%v", utils.CreatePacket("nack"))

			_, err = connection.Write([]byte(myMagicBytes))
		}
		return valToRet, err
	}
}

func (socket *Socket) SendMessage(message string, connection *net.UDPConn) (err error) {

	buffer := make([]byte, constants.BUFF_SIZE)

	var receivedStruct *models.Packet

	myMagicBytes := fmt.Sprintf("%v", utils.CreatePacket(message))
	_, err = connection.Write([]byte(myMagicBytes))

	if err != nil {
		return
	}

	_, _, err = connection.ReadFromUDP(buffer)

	if err != nil {
		return
	}

	err = json.Unmarshal(buffer, &receivedStruct)
	if err != nil {
		return err
	}

	if receivedStruct.Payload != "nack" {
		_, err := connection.Write([]byte(myMagicBytes))
		if err != nil {
			return err
		}
		_, _, err = connection.ReadFromUDP(buffer)

		if err != nil {
			return err
		}
	}

	return nil
}

// connection for client
func (socket *Socket) ClientSocket(udpAddress *net.UDPAddr) (clientConn *net.UDPConn, err error) {
	clientConn, err = net.DialUDP(socket.networkType, nil, udpAddress)
	if err != nil {
		return nil, err
	}

	if err = socket.SendMessage("connect", clientConn); err != nil {
		return nil, err
	}

	return clientConn, nil
}

// connection for server
func (socket *Socket) ServerSocketConnect(udpAddress *net.UDPAddr) (serverConn *net.UDPConn, err error) {

	serverConn, err = net.ListenUDP(socket.networkType, udpAddress)

	if err != nil {
		return nil, err
	}

	return serverConn, err
}

//Create a pointer to udp address;
func (socket *Socket) CreateUpdAddress() (udpAddress *net.UDPAddr, err error) {
	udpAddress, err = net.ResolveUDPAddr(socket.networkType, socket.port)
	if err != nil {
		return nil, err
	}

	return udpAddress, nil
}
