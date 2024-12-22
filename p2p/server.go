package p2p

import (
	"bufio"
	"fmt"
	"local/blockchain-dump/blockchain" // Import the blockchain package
	"net"
	"os"
)

// ========================Global Blockchain Instance========================

// Declare a global blockchain instance
var chain *blockchain.Blockchain

// ========================Peer-to-Peer Server========================

// StartServer initializes the blockchain and starts the P2P server to listen for incoming connections
func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on port", port)

	// Initialize the blockchain instance
	chain = blockchain.InitBlockchain()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

// ========================Message Handling========================

// handleConnection processes incoming messages from a peer and performs actions based on the message type
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		messageJSON, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by peer:", err)
			break
		}

		message, err := DeserializeMessage(messageJSON)
		if err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		switch message.Type {
		case "DATA_HASH":
			dataHash := message.Data.(string)

			// Retrieve data using IPFS
			data, err := blockchain.SimulateIPFSDownload(dataHash)
			if err != nil {
				fmt.Println("Error retrieving data from IPFS:", err)
				continue
			}

			// Process the data using the deterministic XOR algorithm
			output := blockchain.DeterministicAlgorithm(data)

			// Create a transaction and mine a block
			tx := blockchain.Transaction{
				DataHash: dataHash,
				AlgoHash: blockchain.HashData("XOR"),
				Output:   output,
			}

			// Get the latest block and create a new block
			latestBlock := chain.GetLatestBlock()
			newBlock := blockchain.NewBlock([]blockchain.Transaction{tx}, latestBlock.Hash)

			// Add the new block to the blockchain
			chain.AddBlockToChain(newBlock)

			fmt.Println("Mined new block with transaction:", tx)
			BroadcastMessage("NEW_BLOCK", newBlock)

		default:
			fmt.Println("Unknown message type:", message.Type)
		}
	}
}
