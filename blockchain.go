package main

type Blockchain struct {
	chain []*Block
}

func (chain *Blockchain) AddBlock(data string) {
	previousBlock := chain.chain[len(chain.chain)-1]
	newBlock := NewBlock(data, previousBlock.Hash)
	chain.chain = append(chain.chain, newBlock)
}

// We need a starting point, which is why this function is implemented. This is to generate a first time block

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{GenerateFirstBlock()}}
}
