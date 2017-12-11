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

func (cli *CLI) usage() {
	fmt.Println("Usage:")
	fmt.Println(" blockchain_go addblock -data BLOCK_DATA - Adds a block to the blockchain")
	fmt.Println(" format - print all the blocks in the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.format()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Succesfully added a new block...")
}

func (cli *CLI) format() {
	chain := cli.bc.Iterator()
	for {

		block := chain.next()

		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Prev hash: %x\n", block.Previous)
		fmt.Printf("Block Data: %s\n", block.Data)
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

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	formatCmd := flag.NewFlagSet("format", flag.ExitOnError)

	data := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		fmt.Println("x")
		err := addBlockCmd.Parse(os.Args[2:])
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

	if addBlockCmd.Parsed() {
		if *data == "" {

			addBlockCmd.Usage()
			os.Exit(1)
		}
		fmt.Println("y")
		cli.addBlock(*data)
	}

	if formatCmd.Parsed() {
		cli.format()
	}
}
