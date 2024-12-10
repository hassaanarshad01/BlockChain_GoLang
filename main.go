package main

//Module name is BlockchainProject
//go mod init BlockchainProject ------> Command to init module
// Within this module is the following structure:
//	BlockchainProject
//		-blockchain
//			-block.go
//			-PoW.go
//		-main.go


import "BlockchainProject/blockchain"
import "fmt"

func main() {

	//TESTING
	chain := blockchain.InitBlockchain()
	chain.AddBlock("First Block after Genesis Block")
	chain.AddBlock("Second Block after Genesis Block")
	chain.AddBlock("Third Block after Genesis Block")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Block Hash = %x\n", block.PrevHash)
		fmt.Printf("Data = %s\n", block.Data)
		fmt.Println("Hash = %x\n\n", block.Hash)
	}
}