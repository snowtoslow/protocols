package protocol

import (
	"log"
	"net"
	"protocols/utils"
)

func (socket *MyMagicSocket) SecuredSend(connection *net.UDPConn, sharedClient []byte, message string) (err error) {
	myEncryptedByte, err := utils.EncryptionLast(sharedClient, message)
	if err != nil {
		return err
	}

	encrString := myEncryptedByte

	if err = socket.HandleClient(connection, encrString); err != nil {
		return err
	}

	log.Println(myEncryptedByte)

	return nil
}

func (socket *MyMagicSocket) SecuredReceive(connection *net.UDPConn, sharedServer []byte) (err error) {

	validatedPacketToDecrypt, err := socket.HandleServer(connection)
	if err != nil {
		return err
	}

	newStruct, err := utils.CreateStructFromBytes(validatedPacketToDecrypt)
	if err != nil {
		return err
	}

	decryptedBytes, err := utils.DecryptionLast(sharedServer, newStruct.Payload)
	if err != nil {
		return err
	}

	log.Println("DECRYPTED BYTES:", decryptedBytes)

	return nil

}
