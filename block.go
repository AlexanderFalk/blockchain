package main

import "time"

type blockService interface {
	// Add a new node to the list
	RegisterNode(address []byte) bool

	// Create a new block in the blockchain
	NewBlock(proof []byte, previousHash []byte) *Block

	// Returns the last block in the chain
	LastBlock() Block

	// Adding block to the main chain
	AddBlock(data []byte)

	ValidChain(chain *[]Block) bool

	Conflicts() bool
}

type Block struct {
	Index     int64
	Data      []byte
	Timestamp int64
	Hash      []byte
	Previous  []byte
	Proof     []byte
	Nonce     int
}

type Transaction struct {
}

type Token struct {
}

var index int64 = 0

func NewBlock(data string, previousHash []byte) *Block {
	index++
	block := &Block{
		index,
		[]byte(data),
		time.Now().Unix(),
		[]byte{},
		previousHash,
		[]byte{},
		0}

	pow := ShiftBits(block)
	nonce, hash := pow.RunProof()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func GenerateFirstBlock() *Block {
	return NewBlock("Genesis Starting Block", []byte{})
}
