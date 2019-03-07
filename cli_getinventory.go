package main

import (
	"fmt"
	"log"
)

func (cli *CLI) getInventory(address string, nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	//fmt.Printf("%+v\n", UTXOSet.Blockchain)
	defer bc.db.Close()

	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	if len(UTXOs) == 0 {
		fmt.Printf("Empty Inventory of '%s'\n", address)
	} else {
		fmt.Printf("Items in Inventory of '%s'\n", address)

		for count, out := range UTXOs {
			fmt.Printf("Item %d : %s ", count+1, out.Item)
		}
	}

}
