package main

import (
	"fmt"
)

func (cli *CLI) createBlockchain(name, pubKey, gstin, prefix string, nodeID string) {
	CreateBlockchain(name, pubKey, gstin, prefix, nodeID)
	fmt.Println("Done!")
}
