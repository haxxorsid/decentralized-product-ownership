package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// Product struct
type Product struct {
	ID        []byte
	Code      int
	Name      []byte
	Signature []byte
	PubKey    []byte
}

// Serialize returns a serialized Product
func (p Product) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(p)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// Hash returns the hash of the Product
func (p *Product) Hash() []byte {
	var hash [32]byte

	product := *p
	product.ID = []byte{}

	hash = sha256.Sum256(product.Serialize())

	return hash[:]
}

// Deserialize deserializes Product
func DeserializeProduct(data []byte) Product {
	var p Product

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&p)
	if err != nil {
		log.Panic(err)
	}

	return p
}

// Signs a Product
func (p *Product) Sign(privKey ecdsa.PrivateKey, pubKeyHash []byte) {
	pCopy := p.Copy()

	pCopy.Signature = nil
	pCopy.PubKey = pubKeyHash
	pCopy.ID = pCopy.Hash()
	pCopy.PubKey = nil

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, pCopy.ID)
	if err != nil {
		log.Panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)

	p.Signature = signature
}

// String returns a human-readable representation of a product
func (p Product) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Product %x:", p.ID))

	lines = append(lines, fmt.Sprintf("       Name:       %s", p.Name))
	lines = append(lines, fmt.Sprintf("       Code:       %d", p.Code))
	lines = append(lines, fmt.Sprintf("       Signature: %x", p.Signature))
	lines = append(lines, fmt.Sprintf("       PubKey: %x", p.PubKey))
	return strings.Join(lines, "\n")
}

// Copy creates a copy of Product to be used in signing
func (p *Product) Copy() Product {

	pCopy := Product{p.ID, p.Code, p.Name, nil, p.PubKey}

	return pCopy
}

// Verifies a Product
func (p *Product) Verify(pubKeyHash []byte) bool {
	curve := elliptic.P256()
	pCopy := p.Copy()

	pCopy.Signature = nil
	pCopy.PubKey = pubKeyHash
	pCopy.ID = pCopy.Hash()
	pCopy.PubKey = nil

	r := big.Int{}
	s := big.Int{}
	sigLen := len(p.Signature)
	r.SetBytes(p.Signature[:(sigLen / 2)])
	s.SetBytes(p.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(p.PubKey)
	x.SetBytes(p.PubKey[:(keyLen / 2)])
	y.SetBytes(p.PubKey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, pCopy.ID, &r, &s) == false {
		return false
	}

	return true
}

// NewProduct creates a new product
func NewProducts(address string, products []string, nodeID string) ([]*Product, error) {
	var ps []*Product
	var count int

	wallets, err := NewWallets(nodeID)

	if err != nil {
		log.Panic(err)
	}

	ProductCache := ProductCacheSet{NewBlockchain(nodeID)}

	wallet := wallets.GetWallet(address)
	r := ProductCacheSet.Blockchain.GetRole(wallet.PublicKey)

	if r == nil {
		return nil, errors.New("Organisation not found")
	}

	if bytes.Compare(r, []byte("Manufacturer")) != 0 {
		return nil, errors.New("Not authorized to perform this action")
	}

	count = ProductCacheSet.Blockchain.GetNextProductCode(address)

	for index, p := range products[:] {
		product := &Product{nil, count + index, []byte(p), nil, wallet.PublicKey}
		product.ID = product.Hash()
		ps = append(ps, product)
	}

	SignProducts(ps, wallet.PrivateKey, Base58Decode([]byte(address)))

	return ps, nil
}

// SignProducts signs Products
func SignProducts(products []*Product, privKey ecdsa.PrivateKey, pubKeyHash []byte) {
	for _, p := range products[:] {
		p.Sign(privKey, pubKeyHash)
	}
}
