package protocol

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"protocols/constants"
	"protocols/security"
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

func (socket *MyMagicSocket) SendValueToClient(connection *net.UDPConn) (sharedKeyFromServer *big.Int, err error) {

	var clientStructReceived *security.ClintSecuredSendStruct

	buffer := make([]byte, constants.BUFF_SIZE) //4000 in case of shit

	serverSecret, err := security.GeneratePrivateKey()
	if err != nil {
		sharedKeyFromServer = nil
	}

	n, add, err := connection.ReadFromUDP(buffer)
	if err != nil {
		sharedKeyFromServer = nil
	}

	if err := json.Unmarshal(buffer[:n], &clientStructReceived); err != nil {
		sharedKeyFromServer = nil
	}

	computedSharedKeyServer := security.ServerComputes(clientStructReceived, serverSecret)

	log.Println("SHARED VALUE ON SERVER:", computedSharedKeyServer)

	serverSecuredStructToClient, err := security.CreateSecuredStructToSendFromServerToClient(clientStructReceived.FirstPublicNum,
		clientStructReceived.SecondPublicNum, serverSecret)
	if err != nil {
		sharedKeyFromServer = nil
	}

	bytesFromPacket, err := json.Marshal(&serverSecuredStructToClient)
	if err != nil {
		bytesFromPacket = nil
	}

	if _, err := connection.WriteToUDP(bytesFromPacket, add); err != nil {
		sharedKeyFromServer = nil
	}

	n, add, err = connection.ReadFromUDP(buffer)
	if err != nil {
		sharedKeyFromServer = nil
	}

	log.Println("SERVER READ FROM CLIENT SHARED KEY:", string(buffer[:n]), add)
	var sharedClientKey *security.ValuesComputedAfterSend

	if err = json.Unmarshal(buffer[:n], &sharedClientKey); err != nil {
		sharedKeyFromServer = nil
	}

	sharedKeyFromServer = computedSharedKeyServer.Value

	log.Println("CLIENT SHARED ON SERVER:", sharedClientKey.Value)
	log.Printf("IS TRUE: %v", sharedClientKey.Value.Cmp(computedSharedKeyServer.Value)) // if we compare two big ints with Cmp it will return 0 if they are equals;

	if sharedClientKey.Value.Cmp(computedSharedKeyServer.Value) == 0 {
		log.Println("SECURED CONNECTION WAS ESTABLISHED!")
		if _, err := connection.WriteToUDP([]byte("0"), add); err != nil {
			sharedKeyFromServer = nil
		}
	}

	return
}

//func to put in client socket
func (socket *MyMagicSocket) SendPubNumToServer(connection *net.UDPConn) (sharedClient *big.Int, err error) {
	log.Println("Send to server!")

	buffer := make([]byte, constants.BUFF_SIZE) //4000 in case of shit

	var serverSecuredStruct *security.ServerSecuredStruct

	clientSecret, err := security.GeneratePrivateKey()
	if err != nil {
		sharedClient = nil
	}

	publicFirstNum, err := rand.Prime(rand.Reader, 256)
	publicSecondNum, err := rand.Prime(rand.Reader, 256)

	if err != nil {
		sharedClient = nil
	}

	clientStruct, err := security.CreateSecuredStructToSendFromClientToServer(publicFirstNum, publicSecondNum, clientSecret)
	if err != nil {
		sharedClient = nil
	}

	bytesClientStruct, err := json.Marshal(clientStruct)
	if err != nil {
		sharedClient = nil
	}

	if _, err := connection.Write(bytesClientStruct); err != nil {
		sharedClient = nil
	}

	n, _, err := connection.ReadFromUDP(buffer)
	if err != nil {
		sharedClient = nil
	}

	if err := json.Unmarshal(buffer[:n], &serverSecuredStruct); err != nil {
		sharedClient = nil
	}

	sharedValueClient := security.ClientComputes(serverSecuredStruct, publicFirstNum, clientSecret)

	log.Println("SHARED VALUE ON CLIENT:", sharedValueClient)

	sharedBytesStruct, err := json.Marshal(sharedValueClient)
	if err != nil {
		sharedClient = nil
	}

	if _, err := connection.Write(sharedBytesStruct); err != nil {
		sharedClient = nil
	}

	n, _, err = connection.ReadFromUDP(buffer)
	if err != nil {
		sharedClient = nil
	}

	if buffer[:n][0] == 48 {
		sharedClient = sharedValueClient.Value
	}

	return
}

//here is for client
func (socket *MyMagicSocket) SendMessage(connection *net.UDPConn, message string) (err error) {

	myMagicCounter := 5

	buffer := make([]byte, constants.BUFF_SIZE)

	packet := utils.CreatePacket(message)

	bytesFromPacket, err := utils.CreateBytesFromPacket(&packet)

	for i := 0; i < myMagicCounter; i++ {

		if _, err := connection.Write(bytesFromPacket); err != nil {
			return err
		}

		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			return err
		}

		log.Printf("msg: %s, addr: %s", buffer[:n], addr)

		packetAfterRead, err := utils.CreateStructFromBytes(buffer[:n])
		if err != nil {
			return err
		}

		if packetAfterRead.Payload == "acknowledged" {
			break
		}

	}

	return nil
}

//here is for server
func (socket *MyMagicSocket) ReceiveMessage(connection *net.UDPConn) (err error) {

	buffer := make([]byte, constants.BUFF_SIZE)

	n, add, err := connection.ReadFromUDP(buffer)
	if err != nil {
		return err
	}

	//log.Printf("MSG, From addr :%s, %s", buffer[:n], add)

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

/*func (socket *MyMagicSocket)ReceiveRefactored(connection *net.UDPConn) (err error){
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

	log.Println(packet)

	return nil
}


func (socket *MyMagicSocket)SendRefactored(connection *net.UDPConn, message string) error {

	packet := utils.CreatePacket(message)

	bytesFromPacket, err := utils.CreateBytesFromPacket(&packet)

	if err!=nil {
		return err
	}

	if _, err := connection.Write(bytesFromPacket); err != nil{
		return err
	}
	return nil
}*/
