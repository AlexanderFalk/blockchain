package main

import (
	"fmt"
	"strconv"
)

func main() {

	createChain := NewBlockChain()

	createChain.AddBlock("This is first addition to the block")
	createChain.AddBlock("This is the second addition to the block")

	for _, block := range createChain.chain {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %d\n", block.Nonce)
		fmt.Printf("Prev hash: %x\n", block.Previous)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		pow := ShiftBits(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.ValidProof()))
		fmt.Println()
	}

}
