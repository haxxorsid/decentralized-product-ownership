package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -name NAME -publickey KEY -gstin GSTIN -prefix PREFIX - Create blockchains")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  createorg -address ADDRESS -name NAME -publickey KEY -gstin GSTIN -prefix PREFIX -role ROLE - Add a organisation")
	fmt.Println("  listorg - Print all organisations")
	fmt.Println("  inventory -address ADDRESS - Inventory of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet")
	fmt.Println("  printchain - Print all the blocks of the transaction chain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -products PRODUCT -mine - Send PRODUCT from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  addproducts -address ADDRESS -names NAMES -Add products")
	fmt.Println("  listproducts -address ADDRESS -Get products of address")
	fmt.Println("  produceproducts -address ADDRESS -codes CODES -Produce product")
	fmt.Println("  getitemdetails -item ITEM -Get item history")
	fmt.Println("  startnode - Start a node with ID specified in NODE_ID env. var")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	inventoryCmd := flag.NewFlagSet("inventory", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	createOrgCmd := flag.NewFlagSet("createorg", flag.ExitOnError)
	listOrgCmd := flag.NewFlagSet("listorg", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	addProductsCmd := flag.NewFlagSet("addproducts", flag.ExitOnError)
	listProductsCmd := flag.NewFlagSet("listproducts", flag.ExitOnError)
	produceProductsCmd := flag.NewFlagSet("produceproducts", flag.ExitOnError)
	getItemDetailsCmd := flag.NewFlagSet("getitemdetails", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	inventoryAddress := inventoryCmd.String("address", "", "The address to get inventory for")
	createBlockchainName := createBlockchainCmd.String("name", "", "Name of the organisation")
	createBlockchainPublicKey := createBlockchainCmd.String("publickey", "", "PublicKey of the organisation")
	createBlockchainGSTIN := createBlockchainCmd.String("gstin", "", "GSTIN of the organisation")
	createBlockchainPrefix := createBlockchainCmd.String("prefix", "", "Prefix of the organisation")
	produceProductsAddress := produceProductsCmd.String("address", "", "The address to send produced product to")
	createOrgAdminAddr := createOrgCmd.String("address", "", "Address of the admin")
	createOrgName := createOrgCmd.String("name", "", "Name of the organisation")
	createOrgPublicKey := createOrgCmd.String("publickey", "", "PublicKey of the organisation")
	createOrgGSTIN := createOrgCmd.String("gstin", "", "GSTIN of the organisation")
	createOrgPrefix := createOrgCmd.String("prefix", "", "Prefix of the organisation")
	createOrgRole := createOrgCmd.String("role", "", "Role of the organisation")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendProduct := sendCmd.String("products", "", "Item to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	addProductAddress := addProductsCmd.String("address", "", "Product name")
	addProductsName := addProductsCmd.String("names", "", "Product names")
	listProductsAddress := listProductsCmd.String("address", "", "Source Wallet Address")
	cProducts := produceProductsCmd.String("codes", "", "Code of products to produce")
	cItem := getItemDetailsCmd.String("product", "", "Item to find detail of")

	switch os.Args[1] {
	case "inventory":
		err := inventoryCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createorg":
		err := createOrgCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listorg":
		err := listOrgCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addproducts":
		err := addProductsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listproducts":
		err := listProductsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "produceproducts":
		err := produceProductsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getitemdetails":
		err := getItemDetailsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if inventoryCmd.Parsed() {
		if *inventoryAddress == "" {
			inventoryCmd.Usage()
			os.Exit(1)
		}
		cli.getInventory(*inventoryAddress, nodeID)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainName == "" || *createBlockchainPublicKey == "" || *createBlockchainGSTIN == "" || *createBlockchainPrefix == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainName, *createBlockchainPublicKey, *createBlockchainGSTIN, *createBlockchainPrefix, nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}

	if createOrgCmd.Parsed() {
		if *createOrgAdminAddr == "" || *createOrgName == "" || *createOrgPublicKey == "" || *createOrgGSTIN == "" || *createOrgPrefix == "" || *createOrgRole == "" {
			createOrgCmd.Usage()
			os.Exit(1)
		}

		cli.createOrg(*createOrgAdminAddr, *createOrgName, *createOrgPublicKey, *createOrgGSTIN, *createOrgPrefix, *createOrgRole, nodeID)
	}

	if listOrgCmd.Parsed() {
		cli.listOrganisations(nodeID)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printTransactionChain(nodeID)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendProduct == "" {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendProduct, nodeID, *sendMine)
	}

	if addProductsCmd.Parsed() {
		if *addProductAddress == "" || *addProductsName == "" {
			addProductsCmd.Usage()
			os.Exit(1)
		}
		cli.addProducts(*addProductAddress, *addProductsName, nodeID)
	}

	if listProductsCmd.Parsed() {
		if *listProductsAddress == "" {
			listProductsCmd.Usage()
			os.Exit(1)
		}
		cli.listProducts(*listProductsAddress, nodeID)
	}

	if produceProductsCmd.Parsed() {
		if *produceProductsAddress == "" || *cProducts == "" {
			produceProductsCmd.Usage()
			os.Exit(1)
		}

		cli.produceProducts(*produceProductsAddress, *cProducts, nodeID)
	}

	if getItemDetailsCmd.Parsed() {
		if *cItem == "" {
			getItemDetailsCmd.Usage()
			os.Exit(1)
		}

		cli.getItemDetails(*cItem, nodeID)
	}

	if startNodeCmd.Parsed() {
		cli.startNode(nodeID)
	}

}
