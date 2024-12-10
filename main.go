package main

import (
	"BlockchainProject/blockchain"
	"BlockchainProject/p2p"
	"flag"
	"fmt"
)

func main() {
	//usage via: go run main.go --port=3001
	// --port sets the port number, same can be done in dockerfile config
	port := flag.String("port", "3001", "Port to listen on")
	peer := flag.String("peer", "", "Address of a peer to connect to")
	flag.Parse()

	//start the peer up
	go p2p.StartServer(*port)

	//connect to an initial peer if availible.
	if *peer != "" {
		p2p.ConnectToPeer(*peer, "REQUEST_CHAIN")
		p2p.AddPeer(*peer)
	}

	go p2p.CheckPeerHealth()
	chain := blockchain.InitBlockchain()

	chain.AddBlock("First Block after Genesis Block")

	//Broadcast block
	p2p.BroadcastMessage("NEW_BLOCK: First Block after Genesis Block")

	fmt.Println("Blockchain initialized. Listening on port:", *port)
}
