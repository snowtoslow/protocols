package testprot

import (
	"log"
	"net"
	"protocols/constants"
	"protocols/utils"
)

type MyMagicSocket struct {
	network string
	address string
}

func NewMagicSocket(network string, address string) *MyMagicSocket {
	return &MyMagicSocket{
		network: network,
		address: address,
	}
}

func (socket *MyMagicSocket) SendMessage(connection *net.UDPConn, message string) (err error) {

	buffer := make([]byte, constants.BUFF_SIZE)

	packet := utils.CreatePacket(message)

	bytesFromPacket, err := utils.CreateBytesFromPacket(&packet)

	if _, err := connection.Write(bytesFromPacket); err != nil {
		return err
	}

	n, _, _ := connection.ReadFrom(buffer)

	log.Println(string(buffer[:n]))

	return nil
}

func (socket *MyMagicSocket) ReceiveMessage(connection *net.UDPConn) (err error) {

	buffer := make([]byte, constants.BUFF_SIZE)

	n, add, err := connection.ReadFromUDP(buffer)
	if err != nil {
		return err
	}

	log.Printf("MSG, From addr :%s, %s", buffer[:n], add)

	packet, err := utils.CreateStructFromBytes(buffer[:n])
	if err != nil {
		return err
	}

	if utils.ValidatePacket(packet) {
		ackPacket := utils.CreatePacket("acknowledged")
		bytesAcknowledge, err := utils.CreateBytesFromPacket(&ackPacket)
		if err != nil {
			return err
		}

		if _, err := connection.WriteToUDP(bytesAcknowledge, add); err != nil {
			return err
		}

	} else {
		notAck := utils.CreatePacket("not acknowledged")
		bytesNotAcknowledge, err := utils.CreateBytesFromPacket(&notAck)
		if err != nil {
			return err
		}

		if _, err := connection.WriteToUDP(bytesNotAcknowledge, add); err != nil {
			return err
		}
	}

	return nil
}

func (socket *MyMagicSocket) ClientSocket(udpAddress *net.UDPAddr) (connection *net.UDPConn, err error) {
	connection, err = net.DialUDP(socket.network, nil, udpAddress)
	if err != nil {
		connection = nil
	}

	return
}

func (socket *MyMagicSocket) ServerSocket(udpAddress *net.UDPAddr) (connection *net.UDPConn, err error) {
	connection, err = net.ListenUDP(socket.network, udpAddress)
	if err != nil {
		connection = nil
	}
	return
}

func (socket *MyMagicSocket) CreateUdpAddress() (udpAddress *net.UDPAddr, err error) {
	udpAddress, err = net.ResolveUDPAddr(socket.network, socket.address)
	if err != nil {
		udpAddress = nil
	}
	return
}
