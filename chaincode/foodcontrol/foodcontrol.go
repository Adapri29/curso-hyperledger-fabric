package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for controlling the food
type SmartContract struct {
	contractapi.Contract
}

// Food describes basic details of what makes up a food
type Food struct {
	Farmer  string `json:"farmer"`
	Variety string `json:"variety"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, foodId string, farmer string, variety string) error {

	// Validaciones de sintaxis

	// Validaciones de negocio

	food, err := s.Query(ctx, foodId)
	if err == nil && food != nil {
		return fmt.Errorf("food already exists")
	}

	newFood := Food{
		Farmer:  farmer,
		Variety: variety,
	}

	// Transformar el objeto a bytes
	foodAsBytes, err := json.Marshal(newFood)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	// PutState guarda en el libro distribuido
	return ctx.GetStub().PutState(foodId, foodAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, foodId string) (*Food, error) {

	foodAsBytes, err := ctx.GetStub().GetState(foodId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %s", err.Error())
	}

	if foodAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", foodId)
	}

	food := new(Food)

	err = json.Unmarshal(foodAsBytes, food)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s", err.Error())
	}

	return food, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating foodcontrol chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting foodcontrol chaincode: %s", err.Error())
	}
}
