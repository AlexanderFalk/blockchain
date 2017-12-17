package main

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"time"
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
const genesisCoinbaseData = "Borsen 17/Dec/2017 A Treat To The Modern Bank System"

func (chain *Blockchain) MineBlock(transactions []*Transaction) {
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
	newBlock := NewBlock(transactions, prevHash)

	err = chain.database.Update(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blocksBucket))
		err := block.Put(newBlock.Hash, newBlock.SerializeBlock())
		if err != nil {
			log.Panic(err)
		}
		err = block.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		chain.tip = newBlock.Hash

		return nil
	})
}

// Returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Continue statements - A "continue" statement begins the next iteration of the
			// innermost "for" loop at its post statement. The "for" loop must be within the same function.
		Outputs:
			for outIdx, out := range tx.Out {
				// Check if the output was already spent
				// Skipping those that were referenced in inputs.
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// if an output was locked by the same address we're searching unspent transaction outputs for,
				// then this is the output we want.
				if out.CanBeUnockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.In {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.Previous) == 0 {
			break
		}
	}

	return unspentTXs
}

// Takes the transactions and return only outputs
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTXs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTXs {
		for _, out := range tx.Out {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// creates a new blockchain with genesis block
func NewBlockChain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No blockchain has been found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blocksBucket))
		tip = block.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// Creates a new blockchain DB
// We need a starting point, which is why this function is implemented. This is to generate a first time block
func CreateBlockchain(address string) *Blockchain {
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

		if block == nil {
			fmt.Println("No existing block... Creating a new")
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			newGenesis := GenesisBlock(cbtx)

			block, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = block.Put(newGenesis.Hash, newGenesis.SerializeBlock())
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

	return &bc
}

// An iterator initially points at the tip of a blockchain
func (chain *Blockchain) Iterator() *BlockchainIterator {
	it := &BlockchainIterator{chain.tip, chain.database}
	return it
}

func (it *BlockchainIterator) next() *Block {
	var block *Block

	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		decode := b.Get(it.currentHash)
		block = DeserializeBlock(decode)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	it.currentHash = block.Previous

	return block
}
