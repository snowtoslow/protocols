package security

import (
	"crypto/rand"
	"log"
	"math/big"
)

//send to server
type ClintSecuredSendStruct struct {
	FirstPublicNum  *big.Int
	SecondPublicNum *big.Int
	ComputedValue   *big.Int
}

//send to client
type ServerSecuredStruct struct {
	ComputedValue *big.Int
}

type ValuesComputedAfterSend struct {
	Value *big.Int
}

//func to create a struct which will be send to server
func CreateSecuredStructToSendFromClientToServer(firstPublicNum *big.Int, secondPublicNum *big.Int) (clientStruct *ClintSecuredSendStruct, err error) {
	myValue, err := clientComputesValueToSend(firstPublicNum, secondPublicNum)
	if err != nil {
		return nil, err
	}
	log.Println(firstPublicNum, secondPublicNum)
	clientStruct = &ClintSecuredSendStruct{
		FirstPublicNum:  firstPublicNum,
		SecondPublicNum: secondPublicNum,
		ComputedValue:   myValue,
	}

	return clientStruct, nil
}

//create struct with computed value to client after receive client struct
func CreateSecuredStructToSendFromServerToClient(firstPublicNum *big.Int, secondPublicNum *big.Int) (serverSecuredStruct *ServerSecuredStruct, err error) {
	serverSecret, err := GeneratePrivateKey()
	if err != nil {
		serverSecuredStruct = nil
	}

	resultToPower := new(big.Int).Exp(secondPublicNum, serverSecret, nil)

	serverSecuredStruct = &ServerSecuredStruct{
		ComputedValue: new(big.Int).Mod(resultToPower, firstPublicNum),
	}

	return
}

// //shared key in client
func ClientComputes(serverStruct *ServerSecuredStruct, firstGeneratedNum *big.Int) *ValuesComputedAfterSend {
	serverReceivedValues := serverStruct.ComputedValue
	return &ValuesComputedAfterSend{
		Value: serverReceivedValues.Mod(serverReceivedValues, firstGeneratedNum),
	}
}

//shared key in server
func ServerComputes(clientStruct *ClintSecuredSendStruct) *ValuesComputedAfterSend {
	serverReceivedValue := clientStruct.ComputedValue
	return &ValuesComputedAfterSend{
		Value: serverReceivedValue.Mod(serverReceivedValue, clientStruct.FirstPublicNum),
	}
}

func clientComputesValueToSend(firstPublicNum *big.Int, secondPublicNum *big.Int) (computedValue *big.Int, err error) {
	clientSecret, err := GeneratePrivateKey()
	if err != nil {
		clientSecret = nil
	}
	resultToPower := new(big.Int).Exp(secondPublicNum, clientSecret, nil)
	computedValue = new(big.Int).Mod(resultToPower, firstPublicNum)
	return
}

func GeneratePrivateKey() (privateKey *big.Int, err error) {
	privateKey, err = rand.Prime(rand.Reader, 16)
	if err != nil {
		privateKey = nil
	}
	return
}
