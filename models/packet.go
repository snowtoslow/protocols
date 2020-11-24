package models

type Packet struct {
	Payload  string `json:"payload"`
	CheckSum string `json:"checksum"`
}
