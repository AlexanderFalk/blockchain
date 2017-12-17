package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

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
	Index        int64
	Transactions []*Transaction
	Timestamp    int64
	Hash         []byte
	Previous     []byte
	Proof        []byte
	Nonce        int
}

type Token struct {
}

var index int64 = 0

func NewBlock(transactions []*Transaction, previousHash []byte) *Block {
	index++
	block := &Block{
		index,
		transactions,
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

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func GenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// Serializing the block structs
func (b *Block) SerializeBlock() []byte {
	// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	// The zero value for Buffer is an empty buffer ready to use.
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// Deserialize the byte array and turn it into a block
func DeserializeBlock(de []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(de))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
