package main

import (
	"fmt"
)

func (cli *CLI) listOrganisations(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.db.Close()

	if len(bc.tip3) == 0 {
		fmt.Printf("Organisation Blockchain Empty")
	} else {
		bci := bc.Iterator()

		for {
			block := bci.Next3()
			org := block.Organisation
			fmt.Printf("%s\n", org)

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

	}

}
