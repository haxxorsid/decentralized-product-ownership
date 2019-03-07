package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
	"strings"
)

const dbFile = "blockchain_%s.db"
const transactionsBucket = "blocks"
const productsBucket = "products"
const organisationBucket = "organisations"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain implements interactions with a DB
type Blockchain struct {
	tip1 []byte
	tip2 []byte
	tip3 []byte
	db   *bolt.DB
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(name, pubKey, gstin, prefix string, nodeID string) {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(nodeID) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	bc := NewBlockchain(nodeID)
	var tip1 []byte
	var tip2 []byte
	var tip3 []byte

	pub, _ := hex.DecodeString(pubKey)

	org := &Organisation{nil, []byte(name), []byte(strings.ToUpper(gstin)), []byte(prefix), []byte("Admin"), nil, pub, nil}
	org.ID = org.Hash()
	userGenesis := NewBlock(nil, nil, org, nil, bc.GetBestHeight())

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte(utxoBucket))
		if err != nil {
			log.Panic(err)
		}

		b, err := tx.CreateBucket([]byte(transactionsBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), nil)
		if err != nil {
			log.Panic(err)
		}

		b, err = tx.CreateBucket([]byte(productsBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), nil)
		if err != nil {
			log.Panic(err)
		}

		b, err = tx.CreateBucket([]byte(organisationBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(userGenesis.Hash, userGenesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), userGenesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip3 = userGenesis.Hash

		tip1 = nil
		tip2 = nil
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// Reset removes all blockchain data
func (bc *Blockchain) Reset() {

}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionsBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}
		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.tip1 = block.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(nodeID) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip1 []byte
	var tip2 []byte
	var tip3 []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionsBucket))
		tip1 = b.Get([]byte("l"))

		return nil
	})
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(productsBucket))
		tip2 = b.Get([]byte("l"))

		return nil
	})
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(organisationBucket))
		tip3 = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip1, tip2, tip3, db}

	return &bc
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next1()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionsBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionsBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	if len(bc.tip1) != 0 {
		for {
			block := bci.Next1()

			for _, tx := range block.Transactions {
				txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Vout {
					// Was the output spent?
					if spentTXOs[txID] != nil {

						for _, spentOutIdx := range spentTXOs[txID] {
							if spentOutIdx == outIdx {
								continue Outputs
							}
						}

					}

					outs := UTXO[txID]
					outs.Outputs = append(outs.Outputs, out)
					UTXO[txID] = outs
				}

				if tx.IsCoinbase() == false {

					for _, in := range tx.Vin {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}

				}
			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip1, bc.tip2, bc.tip3, bc.db}

	return bci
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
func (bc *Blockchain) GetBlockHashes(typ string) [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		if typ == "t" {
			block := bci.Next1()
		} else if typ == "p" {
			block := bci.Next2()
		} else {
			block := bci.Next3()
		}

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// MineBlock mines a new block with the provided transactions or products
func (bc *Blockchain) MineBlock(transactions []*Transaction, products []*Product, organisation *Organisation, address string) *Block {
	var lastHash []byte
	var lastHeight int

	if transactions != nil {
		for _, tx := range transactions {
			if bc.VerifyTransaction(tx) != true {
				log.Panic("ERROR: Invalid transaction")
			}
		}

		err := bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(transactionsBucket))
			lastHash = b.Get([]byte("l"))
			blockData := b.Get(lastHash)
			block := DeserializeBlock(blockData)

			lastHeight = block.Height
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		newBlock := NewBlock(transactions, nil, nil, lastHash, lastHeight+1)

		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(transactionsBucket))
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			bc.tip1 = newBlock.Hash

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		return newBlock
	} else if products != nil {
		for _, p := range products {
			if p.Verify(Base58Decode([]byte(address))) != true {
				log.Panic("ERROR: Invalid product")
			}
		}

		err := bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(productsBucket))
			lastHash = b.Get([]byte("l"))
			blockData := b.Get(lastHash)
			block := DeserializeBlock(blockData)

			lastHeight = block.Height
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		newBlock := NewBlock(nil, products, nil, lastHash, lastHeight+1)

		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(productsBucket))
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			bc.tip2 = newBlock.Hash

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		return newBlock
	} else {
		for _, p := range products {
			if p.Verify(Base58Decode([]byte(address))) != true {
				log.Panic("ERROR: Invalid organisation")
			}
		}

		err := bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(organisationBucket))
			lastHash = b.Get([]byte("l"))
			blockData := b.Get(lastHash)
			block := DeserializeBlock(blockData)

			lastHeight = block.Height
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		newBlock := NewBlock(nil, nil, organisation, lastHash, lastHeight+1)

		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(organisationBucket))
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			bc.tip3 = newBlock.Hash

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		return newBlock
	}
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// FindProductByCode finds a product by its Code
func (bc *Blockchain) FindProductByCode(address string, Code int) (Product, error) {
	bci := bc.Iterator()

	if len(bc.tip2) != 0 {

		for {
			block := bci.Next2()

			for _, p := range block.Products {
				if Code == p.Code && p.Verify(Base58Decode([]byte(address))) {
					return *p, nil
				}
			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

	}

	return Product{}, errors.New("Product is not found")
}

// FindProductByCode finds a product by its Code
func (bc *Blockchain) GetNextProductCode(address string) int {
	bci := bc.Iterator()

	if len(bc.tip2) != 0 {

		for {
			block := bci.Next2()

			var code int
			found := false

			for _, p := range block.Products {
				if p.Verify(Base58Decode([]byte(address))) {
					found = true
					code = p.Code + 1
				}
			}

			if found {
				return code
			}
			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

	}
	return 1
}

// FindOrganisationByPublicKey finds a organisation by its pubKey
func (bc *Blockchain) FindOrganisationByPublicKey(pubKey []byte) (Organisation, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next3()
		org := block.Organisation
		if bytes.Compare(pubKey, org.PubKey) == 0 {
			return *org, nil
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Organisation{}, errors.New("Organisation not found")
}

// GetRole gives role of organisation by public key
func (bc *Blockchain) GetRole(pubKey []byte) []byte {
	bci := bc.Iterator()

	for {
		block := bci.Next3()
		org := block.Organisation
		if bytes.Compare(pubKey, org.PubKey) == 0 {
			return org.Role
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil
}

// isOrgDuplicate checks if credentials of organisation are duplicate
func (bc *Blockchain) isOrgDuplicate(gstin, prefix, pubKey []byte) bool {
	bci := bc.Iterator()

	for {
		block := bci.Next3()
		org := block.Organisation
		if bytes.Compare(gstin, org.GSTIN) == 0 || bytes.Compare(pubKey, org.PubKey) == 0 || bytes.Compare(prefix, org.Prefix) == 0 {
			return true
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return false
}

// GenerateSGTIN gives SGTIN code by manufacturer pubKey, product code
func (bc *Blockchain) GenerateSGTIN(address string, pubKey []byte, code int) (string, error) {
	bci := bc.Iterator()

	_, err := bc.FindProductByCode(address, code)

	if err != nil {
		return "", err
	}

	org, err := bc.FindOrganisationByPublicKey(pubKey)

	if err != nil {
		return "", err
	}

	if len(bc.tip1) != 0 {
		found := false
		var serial int
		for {
			block := bci.Next1()

			for _, tx := range block.Transactions {
				if tx.IsCoinbase() {
					for _, out := range tx.Vout {
						codeArr := strings.Split(out.Item, ".")
						i, _ := strconv.Atoi(string(codeArr[1]))
						if bytes.Compare([]byte(codeArr[0]), org.Prefix) == 0 && i == code {
							found = true
							serial, _ = strconv.Atoi(string(codeArr[2]))
						}
					}
					if found {
						return string(org.Prefix) + "." + strconv.Itoa(code) + "." + strconv.Itoa(serial+1), nil
					}
				}
			}

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

	}

	return string(org.Prefix) + "." + strconv.Itoa(code) + "." + "1", nil
}
