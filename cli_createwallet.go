package main

import "fmt"

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
	//pubKeyHash := Base58Decode([]byte(address))
	//pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	//fmt.Println(fmt.Sprintf("Inside Lock: %x", pubKeyHash))
	//fmt.Println(fmt.Sprintf("wallet: %s", Base58Encode(pubKeyHash)))
	//wallet := wallets.GetWallet(address)
	//pubKeyHashL := HashPubKey(wallet.PublicKey)
	//fmt.Println(fmt.Sprintf("Inside NewUTXO: %x", pubKeyHashL))
}
