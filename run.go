package main

func main() {
	chain := NewBlockChain()
	defer chain.database.Close()

	cli := CLI{chain}
	cli.Program()
}
