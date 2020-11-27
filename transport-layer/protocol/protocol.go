package protocol

import (
	"encoding/json"
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

	buffer := make([]byte, constants.BUFF_SIZE)
	for {

		n, address, err := connection.ReadFromUDP(buffer)
		if err != nil {
			return "err1", err
		}
		log.Println("Address in Receive:", address, n)
		log.Println("Buffer in Receive: ", string(buffer[:n]))
		if err = json.Unmarshal(buffer[:n], &receivedStruct); err != nil {
			return "err2", err
		}

		log.Println("Structure after unmarshaling in receive:", receivedStruct)

		if utils.ValidatePacket(receivedStruct) {
			b, err := json.Marshal(utils.CreatePacket("connect"))
			if err != nil {
				return "err3", err
			}

			_, err = connection.WriteToUDP(b, address)
			if err != nil {
				return "err4", err
			}

			if receivedStruct.Payload == "connect" {
				//socket.port = address.Port
				log.Println("connection established!")
				return "connection established!", nil
			} else {
				valToRet = receivedStruct.Payload
			}
		} else {
			b, err := json.Marshal(utils.CreatePacket("nack"))
			if err != nil {
				return "err5", err
			}

			_, err = connection.WriteToUDP(b, address)
			if err != nil {
				return valToRet, err
			}
		}

	}
	//return valToRet, err
}

func (socket *Socket) SendMessage(message string, connection *net.UDPConn) (err error) {

	buffer := make([]byte, constants.BUFF_SIZE)

	var receivedStruct *models.Packet

	myMagicPackage := utils.CreatePacket(message)

	b, err := json.Marshal(myMagicPackage)
	if err != nil {
		return err
	}

	_, err = connection.Write(b)
	if err != nil {
		return err
	}

	n, address, err := connection.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	log.Println("ADDRESS in send: ", address)
	log.Println("BUFFER in send: ", buffer)
	err = json.Unmarshal(buffer[:n], &receivedStruct)
	if err != nil {
		return
	}

	if receivedStruct.Payload != "nack" {
		b, err := json.Marshal(myMagicPackage)
		if err != nil {
			return err
		}

		_, err = connection.Write(b)
		if err != nil {
			return err
		}

		_, _, err = connection.ReadFromUDP(buffer)

		if err != nil {
			return err
		}

		log.Println("LAST BUFFER:", string(buffer))
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
