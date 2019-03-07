package main

import (
	"fmt"
)

func (cli *CLI) listProducts(address string, nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	if len(bc.tip2) == 0 {
		fmt.Printf("Product Blockchain Empty")
	} else {
		bci := bc.Iterator()

		for {
			block := bci.Next2()

			for _, p := range block.Products {
				if p.Verify(Base58Decode([]byte(address))) {
					fmt.Printf("Name: %s\n", p.Name)
					fmt.Printf("Code: %d\n", p.Code)
				}
			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}

}
