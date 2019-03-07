package main

const productcacheBucket = "productcache"

// UTXOSet represents UTXO set
type ProductCacheSet struct {
	Blockchain *Blockchain
}

// FindProducts finds products for a public key hash
func (p ProductCacheSet) FindProducts(pubKeyHash []byte) []Product {
	var Products []Product

	return Products
}

// CountProducts returns the number of Products in the ProductCache set
func (p ProductCacheSet) CountProducts() int {
	return 0
}

// Reindex rebuilds the ProductCache set
func (p ProductCacheSet) Reindex() {

}

// Update updates the ProductCache set with product from the Block
// The Block is considered to be the tip2 of a blockchain
func (p ProductCacheSet) Update(block *Block) {

}
