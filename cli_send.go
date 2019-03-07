package main

import (
	"fmt"
	"log"
	"strings"
)

func (cli *CLI) send(from, to string, productstring string, nodeID string, mineNow bool) {
	products := strings.Split(productstring, ",")
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()
	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(wallet, from, to, products, &UTXOSet)

	if mineNow {
		txs := []*Transaction{tx}
		newBlock := bc.MineBlock(txs, nil, nil, "")
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	// newBlock := bc.MineBlock(txs)
	// UTXOSet.Update(newBlock)

	fmt.Println("Success!")
}
