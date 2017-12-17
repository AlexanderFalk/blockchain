// Proof Of Work
package main

// https://en.bitcoin.it/wiki/Proof_of_work

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

type powService interface {
	// Simple Proof of Work
	ShiftBits(lastProof []byte)

	// Validation of proof
	ValidProof(lastProof []byte, proof []byte) bool
}

var maxNonce = math.MaxInt64

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const bitManipulation = 16

func ShiftBits(b *Block) *ProofOfWork {
	bigInt := big.NewInt(1)
	bigInt.Lsh(bigInt, uint(256-bitManipulation))
	pow := &ProofOfWork{b, bigInt}
	return pow
}

// Should be replaced by the ProofOfWork -method . Preparing Data
func (pow *ProofOfWork) GenerateHash(nonce int) []byte {
	timestamp := []byte(strconv.FormatInt(pow.block.Timestamp, 10))
	headers := bytes.Join(
		[][]byte{
			pow.block.Previous,
			pow.block.HashTransactions(),
			timestamp,
			[]byte(strconv.Itoa(nonce))},
		[]byte{})

	return headers
}

func (pow *ProofOfWork) RunProof() (int, []byte) {
	var hashedInteger big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining... \"%s\"\n", pow.block.Transactions)
	for nonce < maxNonce {
		data := pow.GenerateHash(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashedInteger.SetBytes(hash[:])

		if hashedInteger.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) ValidProof() bool {
	var hashedInteger big.Int

	data := pow.GenerateHash(pow.block.Nonce)
	hash := sha256.Sum256(data)
	fmt.Printf("\r%x", hash)
	hashedInteger.SetBytes(hash[:])

	isValid := hashedInteger.Cmp(pow.target) == -1

	return isValid
}
