package main

import (
	"flag"
	"fmt"
	"local/blockchain-dump/blockchain"
	"local/blockchain-dump/p2p"
	"math/rand"
	"time"
)

func main() {
	// Initialize IPFS
	if !initializeIPFS("http://localhost:5001") {
		return
	}

	// Add predefined data to IPFS
	hashes := addPredefinedDataToIPFS([]string{"Option1", "Option2", "Option3"})
	if hashes == nil {
		return
	}

	// Parse command-line arguments
	port, peer := parseFlags()

	// Start the P2P server
	go p2p.StartServer(port)

	// Connect to an initial peer if provided
	if peer != "" {
		connectToInitialPeer(peer)
	}

	// Regularly check peer health
	go p2p.CheckPeerHealth()

	// Initialize the blockchain
	chain := blockchain.InitBlockchain()

	// Trigger deterministic algorithm on random data
	go triggerAlgorithm([]string{"Option1", "Option2", "Option3"})

	// Main application log
	fmt.Println("Genesis block initialized:", chain.Blocks[0])
	fmt.Println("Blockchain initialized. Listening on port:", port)

	// Prevent main from exiting
	select {}
}

// ========================Initialization Functions========================

// initializeIPFS connects to the IPFS client
func initializeIPFS(nodeAddress string) bool {
	err := blockchain.InitializeIPFSClient(nodeAddress)
	if err != nil {
		fmt.Println("Failed to connect to IPFS:", err)
		return false
	}
	fmt.Println("IPFS client initialized successfully.")
	return true
}

// addPredefinedDataToIPFS adds data items to IPFS and returns their hashes
func addPredefinedDataToIPFS(dataItems []string) map[string]string {
	hashes, err := blockchain.AddPredefinedDataToIPFS(dataItems)
	if err != nil {
		fmt.Println("Error adding predefined data to IPFS:", err)
		return nil
	}

	for data, hash := range hashes {
		fmt.Printf("Data: %s, Hash: %s\n", data, hash)
	}

	return hashes
}

// parseFlags parses command-line arguments for the port and peer
func parseFlags() (string, string) {
	port := flag.String("port", "3001", "Port to listen on")
	peer := flag.String("peer", "", "Address of a peer to connect to")
	flag.Parse()
	return *port, *peer
}

// connectToInitialPeer connects to the specified peer
func connectToInitialPeer(peer string) {
	p2p.ConnectToPeer(peer, "REQUEST_CHAIN", nil)
	p2p.AddPeer(peer)
}

// ========================Application Logic========================

// triggerAlgorithm selects random data, computes its hash, and sends it to a random peer
func triggerAlgorithm(dataOptions []string) {
	rand.Seed(time.Now().UnixNano())
	data := dataOptions[rand.Intn(len(dataOptions))]
	dataHash := blockchain.HashData(data)
	p2p.SendDataHashToRandomPeer(dataHash)
}
