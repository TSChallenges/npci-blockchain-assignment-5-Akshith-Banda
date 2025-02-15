package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Car struct {
	ID        string `json:"id"`
	Company   string `json:"company"`
	Model     string `json:"model"`
	Year      int    `json:"year"`
	ChassisNo string `json:"chassis_no"`
	Color     string `json:"color"`
	Owner     string `json:"owner"`
	ForSale   bool   `json:"for_sale"`
	Price     int    `json:"price,omitempty"`
	Bids      []Bid  `json:"bids,omitempty"`
}

type Bid struct {
	Bidder string `json:"bidder"`
	Amount int    `json:"amount"`
}

type CarContract struct {
	contractapi.Contract
}

func (c *CarContract) RegisterCar(ctx contractapi.TransactionContextInterface, id, company, model, chassisNo, color, owner string, year int) error {
	car := Car{
		ID:        id,
		Company:   company,
		Model:     model,
		Year:      year,
		ChassisNo: chassisNo,
		Color:     color,
		Owner:     owner,
		ForSale:   false,
	}

	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutPrivateData("CarCollection", id, carJSON)
	if err != nil {
		return err
	}

	eventData := map[string]interface{}{
		"id":         id,
		"company":    company,
		"model":      model,
		"year":       year,
		"chassis_no": chassisNo,
		"color":      color,
	}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	return ctx.GetStub().SetEvent("CarRegistered", eventJSON)
}

func (c *CarContract) ListCarForSale(ctx contractapi.TransactionContextInterface, id string, price int) error {
	carJSON, err := ctx.GetStub().GetPrivateData("CarCollection", id)
	if err != nil || carJSON == nil {
		return fmt.Errorf("Car not found in private collection")
	}

	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return err
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}

	if car.Owner != clientID {
		return fmt.Errorf("Only the owner can list the car for sale")
	}

	car.ForSale = true
	car.Price = price

	publicCar := struct {
		ID        string `json:"id"`
		Company   string `json:"company"`
		Model     string `json:"model"`
		Year      int    `json:"year"`
		ChassisNo string `json:"chassis_no"`
		Color     string `json:"color"`
		ForSale   bool   `json:"for_sale"`
		Price     int    `json:"price,omitempty"`
	}{
		ID:        car.ID,
		Company:   car.Company,
		Model:     car.Model,
		Year:      car.Year,
		ChassisNo: car.ChassisNo,
		Color:     car.Color,
		ForSale:   car.ForSale,
		Price:     car.Price,
	}

	updatedCarJSON, err := json.Marshal(publicCar)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedCarJSON)
}

func (c *CarContract) TransferOwnership(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	carJSON, err := ctx.GetStub().GetPrivateData("CarCollection", id)
	if err != nil || carJSON == nil {
		return fmt.Errorf("Car not found in private collection")
	}

	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return err
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}

	if car.Owner != clientID {
		return fmt.Errorf("Only the owner can transfer ownership")
	}

	car.Owner = newOwner
	car.ForSale = false
	car.Price = 0
	car.Bids = nil

	updatedCarJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutPrivateData("CarCollection", id, updatedCarJSON)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(id)
}

func main() {
	carChaincode, err := contractapi.NewChaincode(&CarContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating CarContract chaincode: %s", err))
	}

	if err := carChaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting CarContract chaincode: %s", err))
	}
}
