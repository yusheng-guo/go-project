package main

import "fmt"

func main() {
	blockchain := NewBlockchain(2)

	fmt.Println("add a record")
	blockchain.addBlock("Alice", "Bob", 5)
	fmt.Println("success")
	fmt.Println("add a record")
	blockchain.addBlock("John", "Bob", 2)
	fmt.Println("success")

	fmt.Println(blockchain.isValid())

	for _, block := range blockchain.chain[1:] {
		fmt.Println(block.data)
	}
}
