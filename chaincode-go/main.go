package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Book struct {
	ID        string   `json:"ID"`
	Title     string   `json:"title"`
	NumAvail  int      `json:"numAvail"`
	NumBorrow int      `json:"numBorrow"`
	Borrowers []string `json:"borrowers"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Book{
		{ID: "book1", Title: "blue", NumAvail: 6, NumBorrow: 1, Borrowers: []string{"human1"}},
		{ID: "book2", Title: "red", NumAvail: 5, NumBorrow: 2, Borrowers: []string{"human1", "human2"}},
		{ID: "book3", Title: "green", NumAvail: 0, NumBorrow: 3, Borrowers: []string{"human1", "human2", "human3"}},
		{ID: "book4", Title: "yellow", NumAvail: 3, NumBorrow: 2, Borrowers: []string{"human1", "human4"}},
		{ID: "book5", Title: "black", NumAvail: 2, NumBorrow: 1, Borrowers: []string{"human5"}},
		{ID: "book6", Title: "white", NumAvail: 3, NumBorrow: 0, Borrowers: []string{}},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateBook(ctx contractapi.TransactionContextInterface, id string, title string, numAvail int, numBorrow int, borrowers []string) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the book %s already exists", id)
	}
	if numAvail < 0 || numBorrow < 0 {
		return fmt.Errorf("wrong number is given")
	}

	asset := Book{
		ID:        id,
		Title:     title,
		NumAvail:  numAvail,
		NumBorrow: numBorrow,
		Borrowers: borrowers,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadBook(ctx contractapi.TransactionContextInterface, id string) (*Book, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the book %s does not exist", id)
	}

	var asset Book
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateBook(ctx contractapi.TransactionContextInterface, id string, title string, numAvail int, numBorrow int, borrowers []string) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the book %s does not exist", id)
	}

	// Check whether the length of borrowers array matches numBorrow
	if len(borrowers) != numBorrow {
		return fmt.Errorf("number of borrowers does not match")
	}

	// overwriting original asset with new asset
	asset := Book{
		ID:        id,
		Title:     title,
		NumAvail:  numAvail,
		NumBorrow: numBorrow,
		Borrowers: borrowers,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteBook(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.BookExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the book %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) BookExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) PurchaseBook(ctx contractapi.TransactionContextInterface, id string) error {
	book, err := s.ReadBook(ctx, id)
	if err != nil {
		return err
	}
	if book.NumAvail == 0 {
		return fmt.Errorf("the book %s is not available", id)
	}

	book.NumAvail = book.NumAvail - 1
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) BorrowBook(ctx contractapi.TransactionContextInterface, id string, borrower string) error {
	book, err := s.ReadBook(ctx, id)
	if err != nil {
		return err
	}
	if book.NumAvail == 0 {
		return fmt.Errorf("the book %s is not available", id)
	}

	book.NumAvail = book.NumAvail - 1
	book.NumBorrow = book.NumBorrow + 1
	book.Borrowers = append(book.Borrowers, borrower)
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) ReturnBook(ctx contractapi.TransactionContextInterface, id string, borrower string) error {
	book, err := s.ReadBook(ctx, id)
	if err != nil {
		return err
	}

	found := Find(book.Borrowers, borrower)
	if found == -1 {
		return fmt.Errorf("the book %s does not have borrower %s", id, borrower)
	}

	book.NumAvail = book.NumAvail + 1
	book.NumBorrow = book.NumBorrow - 1
	book.Borrowers = append(book.Borrowers[:found], book.Borrowers[found+1:]...)
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Book, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Book
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Book
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
