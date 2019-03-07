package main

import (
	"fmt"
	"log"
	"strings"
)

func (cli *CLI) addProducts(address string, name string, nodeID string) {
	products := strings.Split(name, ",")
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	
	productArr, err := NewProducts(address, products, nodeID)
	if err != nil {
		fmt.Println(err)
	} else {
		bc := NewBlockchain(nodeID)
		defer bc.db.Close()
		bc.MineBlock(nil, productArr, nil, address)
		fmt.Println("Success!")
	}
}
