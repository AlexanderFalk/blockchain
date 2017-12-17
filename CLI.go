package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) createBlockchain(address string) {
	chain := CreateBlockchain(address)
	chain.database.Close()
	fmt.Println("Created a new blockchain!")
}

func (cli *CLI) usage() {
	fmt.Println("Usage:")
	fmt.Println(" createblockchain -address ADDRESS - Creates a blockchain and send genesis block reward to ADDRESS")
	fmt.Println(" format - print all the blocks in the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.format()
	}
}

// The account balance is the sum of values of all unspent transactions outputs locked by the account address
func (cli *CLI) getBalance(address string) {
	chain := NewBlockChain(address)
	defer chain.database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) format() {
	chain := cli.bc.Iterator()
	for {

		block := chain.next()

		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Prev hash: %x\n", block.Previous)
		fmt.Printf("Block Transctions: %s\n", block.Transactions)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		pow := ShiftBits(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.ValidProof()))
		fmt.Println()

		if len(block.Previous) == 0 {
			break
		}
	}
}

func (cli *CLI) Program() {
	cli.validateArgs()

	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	formatCmd := flag.NewFlagSet("format", flag.ExitOnError)

	//data := addBlockCmd.String("data", "", "Block data")
	createBlockchainAddress := createBlockchainCmd.String("address", "value", "The address to send the reward to")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "format":
		err := formatCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.usage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {

			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if formatCmd.Parsed() {
		cli.format()
	}
}
