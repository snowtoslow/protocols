package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"protocols/models"
	"strings"
)

func ValidatePacket(packet *models.Packet) bool {
	return packet.CheckSum == getMD5HAsh(packet.Payload)
}

func CreatePacket(message string) models.Packet {
	return models.Packet{
		Payload:  strings.TrimRight(message, "\n"),
		CheckSum: getMD5HAsh(strings.TrimRight(message, "\n")),
	}
}

func getMD5HAsh(message string) string {
	myMagicHasher := md5.New()
	myMagicHasher.Write([]byte(message))
	return hex.EncodeToString(myMagicHasher.Sum(nil))
}

func CreateBytesFromPacket(packet *models.Packet) (bytesFromPacket []byte, err error) {
	bytesFromPacket, err = json.Marshal(packet)
	if err != nil {
		bytesFromPacket = nil
	}
	return
}

func CreateStructFromBytes(bytesFromPacket []byte) (packetFromBytes *models.Packet, err error) {
	if err = json.Unmarshal(bytesFromPacket, &packetFromBytes); err != nil {
		packetFromBytes = nil
	}
	return
}

/*b, _ := json.Marshal(utils.CreatePacket("connect"))
log.Println(string(b))
var receivedStruct *models.Packet
// Convert bytes to string.
err := json.Unmarshal(b, &receivedStruct)

if err != nil {
	log.Println(err)
}

log.Println("Struct:",receivedStruct)*/
