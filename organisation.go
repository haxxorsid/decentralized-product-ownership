package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// Organisation struct
type Organisation struct {
	ID          []byte
	Name        []byte
	GSTIN       []byte
	Prefix      []byte
	Role        []byte
	Signature   []byte
	PubKey      []byte
	AdminPubKey []byte
}

// Serialize returns a serialized Organisation
func (org Organisation) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(org)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// Hash returns the hash of the Organisation
func (org *Organisation) Hash() []byte {
	var hash [32]byte

	o := *org
	o.ID = []byte{}

	hash = sha256.Sum256(o.Serialize())

	return hash[:]
}

// Deserialize deserializes Organisation
func DeserializeOrganisation(data []byte) Organisation {
	var org Organisation

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&org)
	if err != nil {
		log.Panic(err)
	}

	return org
}

// Signs a Organisation
func (org *Organisation) Sign(privKey ecdsa.PrivateKey, pubKeyHash []byte) {
	oCopy := org.Copy()

	oCopy.Signature = nil
	oCopy.AdminPubKey = pubKeyHash
	oCopy.ID = oCopy.Hash()
	oCopy.AdminPubKey = nil

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, oCopy.ID)
	if err != nil {
		log.Panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)

	org.Signature = signature
}

// String returns a human-readable representation of a organisation
func (org Organisation) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Organisation %x:", org.ID))

	lines = append(lines, fmt.Sprintf("       Name:       %s", org.Name))
	lines = append(lines, fmt.Sprintf("       GSTIN:       %s", org.GSTIN))
	lines = append(lines, fmt.Sprintf("       Prefix:       %s", org.Prefix))
	lines = append(lines, fmt.Sprintf("       Role:       %s", org.Role))
	lines = append(lines, fmt.Sprintf("       Signature: %x", org.Signature))
	lines = append(lines, fmt.Sprintf("       PubKey: %x", org.PubKey))
	lines = append(lines, fmt.Sprintf("       AdminPubKey: %x", org.AdminPubKey))
	return strings.Join(lines, "\n")
}

// Copy creates a copy of Organisation to be used in signing
func (org *Organisation) Copy() Organisation {

	oCopy := Organisation{org.ID, org.Name, org.GSTIN, org.Prefix, org.Role, nil, org.PubKey, org.AdminPubKey}

	return oCopy
}

// Verifies a Organisation
func (org *Organisation) Verify(pubKeyHash []byte) bool {
	curve := elliptic.P256()
	oCopy := org.Copy()

	oCopy.Signature = nil
	oCopy.AdminPubKey = pubKeyHash
	oCopy.ID = oCopy.Hash()
	oCopy.AdminPubKey = nil

	r := big.Int{}
	s := big.Int{}
	sigLen := len(org.Signature)
	r.SetBytes(org.Signature[:(sigLen / 2)])
	s.SetBytes(org.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(org.AdminPubKey)
	x.SetBytes(org.AdminPubKey[:(keyLen / 2)])
	y.SetBytes(org.AdminPubKey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, oCopy.ID, &r, &s) == false {
		return false
	}

	return true
}

// NewOrganisation creates a new organisation
func NewOrganisation(address, name, pubKey, gstin, prefix, role string, OrganisationCache *OrganisationCacheSet) (*Organisation, error) {

	wallets, err := NewWallets()

	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(address)
	r := OrganisationCache.Blockchain.GetRole(wallet.PublicKey)

	if r == nil {
		return &Organisation{}, errors.New("Organisation not found")
	}

	if bytes.Compare(r, []byte("Admin")) != 0 {
		return &Organisation{}, errors.New("Not authorized to perform this action")
	}

	gst := []byte(strings.ToUpper(gstin))
	pre := []byte(prefix)
	pub, _ := hex.DecodeString(pubKey)

	if OrganisationCache.Blockchain.isOrgDuplicate(gst, pre, pub) {
		return &Organisation{}, errors.New("Duplicate data cannot be added")
	}

	org := &Organisation{nil, []byte(name), gst, pre, []byte(role), nil, pub, wallet.PublicKey}
	org.ID = org.Hash()

	org.Sign(wallet.PrivateKey, Base58Decode([]byte(address)))

	return org, nil
}
