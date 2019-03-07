package main

import (
	"fmt"
	"log"
)

func (cli *CLI) createOrg(address, name, publicKey, gstin, prefix, role, nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain(nodeID)
	OrganisationCache := OrganisationCacheSet{bc}
	defer bc.db.Close()
	org, err := NewOrganisation(address, name, publicKey, gstin, prefix, role, &OrganisationCache)
	if err != nil {
		fmt.Println(err)
	} else {
		bc.MineBlock(nil, nil, org, address)
		fmt.Println("Success!")
	}
}
