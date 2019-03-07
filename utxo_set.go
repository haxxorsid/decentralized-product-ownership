package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

// UTXOSet represents UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, products []string) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	found := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		count := 0

		for k, v := c.First(); k != nil; k, v = c.Next() {
			count = count + 1
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {

				for _, product := range products {
					if out.IsLockedWithKey(pubkeyHash) && found < len(products) && product == out.Item {
						found++
						unspentOutputs[txID] = append(unspentOutputs[txID], out.Index)
					}
				}

			}

		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return found, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				//fmt.Printf("BEFORE IF %+v\n", out)
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
					//fmt.Printf("%+v\n", out)

				}
			}

		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip1 of a blockchain
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {

				for _, vin := range tx.Vin {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)

					for _, out := range outs.Outputs {

						if out.Index != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := b.Put(vin.Txid, updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}

			}

			newOutputs := TXOutputs{}

			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			fmt.Printf("ID: %x \n", tx.ID)
			fmt.Printf("Serial: %x \n", newOutputs.Serialize())
			err := b.Put(tx.ID, newOutputs.Serialize())

			if err != nil {
				fmt.Printf("INSIDE \n")
				log.Panic(err)
			}
			fmt.Printf("TEST \n")
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
