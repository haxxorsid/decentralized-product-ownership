package main

import (
	"fmt"
	"log"
)

func (cli *CLI) listAddresses(nodeID string) {
	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Printf("Address: %s\n", address)
		wallet := wallets.GetWallet(address)
		fmt.Printf(fmt.Sprintf("PubKey: %x \n", wallet.PublicKey))
	}

}
