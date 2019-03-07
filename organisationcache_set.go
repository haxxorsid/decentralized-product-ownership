package main

const organisationcacheBucket = "productcache"

// UTXOSet represents UTXO set
type OrganisationCacheSet struct {
	Blockchain *Blockchain
}

// FindOrganisation finds organisation for a public key
func (p OrganisationCacheSet) FindOrganisation(pubKeyHash []byte) []Product {
	var Products []Product

	return Products
}

// Reindex rebuilds the OrganisationCache set
func (p OrganisationCacheSet) Reindex() {

}

// Update updates the OrganisationCache set with organisation from the Block
// The Block is considered to be the tip3 of a blockchain
func (p OrganisationCacheSet) Update(block *Block) {

}
