package protocol

import (
	"log"
	"net"
	"protocols/utils"
)

func (socket *MyMagicSocket) SecuredSend(connection *net.UDPConn, sharedClient []byte, message string) (err error) {

	myEncryptedByte, err := utils.Encrypt(sharedClient, []byte(message))
	if err != nil {
		return err
	}

	encrString := string(myEncryptedByte)

	if err = socket.HandleClient(connection, encrString); err != nil {
		return err
	}

	return nil
}

func (socket *MyMagicSocket) SecuredReceive(connection *net.UDPConn, sharedServer []byte) (err error) {
	bytes, addr, err := socket.CreateStruct(connection)
	if err != nil {
		return err
	}

	validatedBytes, err := socket.ValidateMsg(connection, bytes, addr)

	decryptedBytes, err := utils.Decrypt(sharedServer, validatedBytes)
	if err != nil {
		return err
	}

	log.Println("DECRYPTED BYTES:", string(decryptedBytes))

	return nil
}
