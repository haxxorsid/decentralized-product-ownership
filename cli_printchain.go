package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) printTransactionChain(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	if len(bc.tip1) == 0 {
		fmt.Printf("Transaction Blockchain Empty")
	} else {
		bci := bc.Iterator()

		for {
			block := bci.Next1()

			fmt.Printf("============ Block %x ============\n", block.Hash)
			fmt.Printf("Height: %d\n", block.Height)
			fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
			pow := NewProofOfWork(block)
			fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))

			for _, tx := range block.Transactions {
				fmt.Println(tx)
			}

			fmt.Printf("\n\n")

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}
}
