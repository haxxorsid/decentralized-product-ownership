package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	Products      []*Product
	Organisation  *Organisation
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, products []*Product, organisation *Organisation, prevBlockHash []byte, height int) *Block {
	var block *Block
	if transactions != nil {
		block = &Block{time.Now().Unix(), transactions, nil, nil, prevBlockHash, []byte{}, 0, height}
	} else if products != nil {
		block = &Block{time.Now().Unix(), nil, products, nil, prevBlockHash, []byte{}, 0, height}
	} else {
		block = &Block{time.Now().Unix(), nil, nil, organisation, prevBlockHash, []byte{}, 0, height}
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// HashTransactionsOrProducts returns a hash of the transactions or products in the block
func (b *Block) HashTransactionsOrProducts() []byte {
	if b.Transactions != nil {
		var transactions [][]byte

		for _, tx := range b.Transactions {
			transactions = append(transactions, tx.Serialize())
		}
		mTree := NewMerkleTree(transactions)

		return mTree.RootNode.Data
	} else if b.Products != nil {
		var products [][]byte

		for _, p := range b.Products {
			products = append(products, p.Serialize())
		}
		mTree := NewMerkleTree(products)

		return mTree.RootNode.Data
	} else {
		var organisations [][]byte
		organisations = append(organisations, b.Organisation.Serialize())
		mTree := NewMerkleTree(organisations)
		return mTree.RootNode.Data
	}
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
