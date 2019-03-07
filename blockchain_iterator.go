package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash1 []byte
	currentHash2 []byte
	currentHash3 []byte
	db           *bolt.DB
}

// Next1 returns next block starting from the tip1
func (i *BlockchainIterator) Next1() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionsBucket))
		encodedBlock := b.Get(i.currentHash1)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash1 = block.PrevBlockHash

	return block
}

// Next2 returns next block starting from the tip2
func (i *BlockchainIterator) Next2() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(productsBucket))
		encodedBlock := b.Get(i.currentHash2)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash2 = block.PrevBlockHash

	return block
}

// Next3 returns next block starting from the tip3
func (i *BlockchainIterator) Next3() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(organisationBucket))
		encodedBlock := b.Get(i.currentHash3)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash3 = block.PrevBlockHash

	return block
}
