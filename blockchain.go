package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type Blockchain struct {
	tip      []byte
	database *bolt.DB
}

// Tx represents a read-only or read/write transaction on the database.
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

const dbFile = "blockchain.db"

const blocksBucket = "blocks"

func (chain *Blockchain) AddBlock(data string) {
	var prevHash []byte

	// Get the last block hash from the DB to use it to mine a new block hash.
	err := chain.database.View(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blocksBucket))
		prevHash = block.Get([]byte("1"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// After mining a new block, we save its serialized representation into the DB and update the l key, which now stores the new block’s hash.
	newBlock := NewBlock(data, prevHash)

	err = chain.database.Update(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blocksBucket))
		err := block.Put(newBlock.Hash, newBlock.SerializeBlock())
		if err != nil {
			log.Panic(err)
		}
		err = chain.database.Update(func(tx *bolt.Tx) error {
			block := tx.Bucket([]byte(blocksBucket))
			err := block.Put(newBlock.Hash, newBlock.SerializeBlock())
			if err != nil {
				log.Panic(err)
			}
			err = block.Put([]byte("1", newBlock.Hash))
			chain.tip = newBlock.Hash

			return nil
		})
	})

	chain.chain = append(chain.chain, newBlock)
}

// We need a starting point, which is why this function is implemented. This is to generate a first time block

func NewBlockChain() *Blockchain {
	var tip []byte
	// Initializes the reference to the database.
	// It's responsible for creating the database if it doesn't exist, obtaining an exclusive lock on the file,
	// reading the meta pages, & memory-mapping the file.
	// 0600 - FILEMODE - Permission : Owner can read and write
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Panic(err)
	}
	// tx represents the internal transaction identifier
	err = db.Update(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blocksBucket))

		if block == null {
			fmt.Println("No existing block... Creating a new")
			newGenesis := GenesisBlock()

			block, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = block.Put(newGenesis.Hash, newGenesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// If it exists, we read the l key from it; if it doesn’t exist,
			// we generate the genesis block, create the bucket, save the block into it,
			// and update the l key storing the last block hash of the chain.
			err = block.Put([]byte("1"), newGenesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = newGenesis.Hash
		} else {
			tip = block.Get([]byte("1"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &Blockchain{[]*Block{GenesisBlock()}}
}

// An iterator initially points at the tip of a blockchain
func (chain *Blockchain) Iterator() *BlockchainIterator {
	it := &BlockchainIterator{chain.tip, chain.database}
	return it
}

func (it *BlockchainIterator) next() *Block {
	var block *Block

	err := it.db.View(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blocksBucket))
		decode := block.Get(it.currentHash)
		block = DeserializeBlock(decode)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return block
}
