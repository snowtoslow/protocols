package security

import (
	"crypto/rand"
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

type SharedKeyStruct struct {
	SharedKey *big.Int
}

//func to create a struct which will be send to server
func CreateSecuredStructToSendFromClientToServer(firstPublicNum *big.Int, secondPublicNum *big.Int, clientSecret *big.Int) (clientStruct *ClintSecuredSendStruct, err error) {
	myValue, err := clientComputesValueToSend(firstPublicNum, secondPublicNum, clientSecret)
	if err != nil {
		return nil, err
	}

	clientStruct = &ClintSecuredSendStruct{
		FirstPublicNum:  firstPublicNum,
		SecondPublicNum: secondPublicNum,
		ComputedValue:   myValue,
	}
	//here is ok
	return clientStruct, nil
}

//create struct with computed value to client after receive client struct
func CreateSecuredStructToSendFromServerToClient(firstPublicNum *big.Int, secondPublicNum *big.Int, serverSecret *big.Int) (serverSecuredStruct *ServerSecuredStruct, err error) {
	resultToPower := new(big.Int).Exp(secondPublicNum, serverSecret, nil)
	valueToSend := new(big.Int).Mod(resultToPower, firstPublicNum)
	serverSecuredStruct = &ServerSecuredStruct{
		ComputedValue: valueToSend,
	}
	//here also is ok

	return
}

// //shared key in client
func ClientComputes(serverStruct *ServerSecuredStruct, firstGeneratedNum *big.Int, secretClient *big.Int) *ValuesComputedAfterSend {
	//serverValue to power of secretClient mod first
	serverReceivedValues := new(big.Int).Exp(serverStruct.ComputedValue, secretClient, nil)
	valuesComputed := new(big.Int).Mod(serverReceivedValues, firstGeneratedNum)
	return &ValuesComputedAfterSend{
		Value: valuesComputed,
	}
}

//shared key in server
func ServerComputes(clientStruct *ClintSecuredSendStruct, secretServer *big.Int) *ValuesComputedAfterSend {
	//client value send to server to power of secretServer mod first
	myVal1 := new(big.Int).Exp(clientStruct.ComputedValue, secretServer, nil)
	valuesComputed := new(big.Int).Mod(myVal1, clientStruct.FirstPublicNum)
	return &ValuesComputedAfterSend{
		Value: valuesComputed,
	}
}

func clientComputesValueToSend(firstPublicNum *big.Int, secondPublicNum *big.Int, clientSecret *big.Int) (computedValue *big.Int, err error) {
	resultToPower := new(big.Int).Exp(secondPublicNum, clientSecret, nil)
	computedValue = new(big.Int).Mod(resultToPower, firstPublicNum)
	return
}

func GeneratePrivateKey() (privateKey *big.Int, err error) {
	privateKey, err = rand.Prime(rand.Reader, 8)
	if err != nil {
		privateKey = nil
	}
	return
}
