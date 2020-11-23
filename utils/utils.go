package utils

import (
	"crypto/md5"
	"encoding/hex"
	"protocols/models"
)

func ValidatePacket(packet *models.Packet) bool {
	return packet.CheckSum == getMD5HAsh(packet.Payload)
}

func CreatePacket(message string) models.Packet {
	return models.Packet{
		Payload:  message,
		CheckSum: getMD5HAsh(message),
	}
}

func getMD5HAsh(message string) string {
	myMagicHasher := md5.New()
	myMagicHasher.Write([]byte(message))
	return hex.EncodeToString(myMagicHasher.Sum(nil))
}
