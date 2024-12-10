package main

import (
	"BlockchainProject/blockchain"
	"BlockchainProject/p2p"
	"fmt"
	"os"
)

func main() {
	
	//Blockchain initialization
	chain := blockchain.InitBlockchain()
	
	//Initialize server port 3000
	go p2p.StartServer("3000") 

	//Block adding simulation
	chain.AddBlock("First Block after Genesis Block")
	chain.AddBlock("Second Block after Genesis Block")

	//sending blockchain to peer
	if len(os.Args) > 1 && os.Args[1] == "connect" {
		p2p.ConnectToPeer("localhost:3000", "REQUEST_CHAIN")
	}

	//printing blockchain
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n\n", block.Hash)
	}
}
