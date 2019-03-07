package main

import (
	"fmt"
)

func (cli *CLI) getItemDetails(item string, nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	if len(bc.tip1) == 0 {
		fmt.Printf("Transaction Blockchain Empty")
	} else {
		bci := bc.Iterator()

		for {
			block := bci.Next1()

			//fmt.Printf("============ Block %x ============\n", block.Hash)
			//fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
			//pow := NewProofOfWork(block)
			//fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))

			for _, tx := range block.Transactions {

				for _, out := range tx.Vout {
					//fmt.Println("out" + out.Item)
					//fmt.Println(item)
					if out.Item == item {
						fmt.Println(fmt.Sprintf("%s", Base58Encode(out.PubKeyHash)))
					}
				}

			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}
}
