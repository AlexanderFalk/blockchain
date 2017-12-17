package main

/*
	When a miner starts mining a block, it adds a coinbase transaction to it.
	A coinbase transaction is a special type of transactions. It doesn't require
	a previously existing output.
	It creates outputs out of nowhere.
*/

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

// The amount of reward from mining. Static in our case
const subsidy = 10

type Transaction struct {
	ID  []byte
	In  []TXInput
	Out []TXOutput
}

type TXOutput struct {
	Value  int
	PubKey string
}

// An input references a previous output
// If all data is correct, the output can be unlocked and the value can be used
// to generate new outputs.
type TXInput struct {
	Txid      []byte // Stores the ID of the output reference
	Out       int    // Stores an index of an output in the transaction
	Signature string // Data to be used in an output's PubKey
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.In) == 1 && len(tx.In[0].Txid) == 0 && tx.In[0].Out == -1
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func (tx Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.Signature == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.PubKey == unlockingData
}
