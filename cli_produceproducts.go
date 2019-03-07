package main

import (
	"fmt"
	"log"
	"strings"
)

func (cli *CLI) produceProducts(address string, productcodes string, nodeID string) {
	codes := strings.Split(productcodes, ",")
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()
	cbTx, err := NewCoinbaseTX(&UTXOSet, address, codes, nodeID)
	if err != nil {
		fmt.Println(err)
	} else {
		txs := []*Transaction{cbTx}
		newBlock := bc.MineBlock(txs, nil, nil, "")
		UTXOSet.Update(newBlock)
		fmt.Println("Success!")
	}
}
