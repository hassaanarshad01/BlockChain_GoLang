package p2p

import (
	"BlockchainProject/blockchain"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// ========================Peer Management========================

var peers []string // List of connected peers
var mu sync.Mutex  // Mutex for thread-safe access to the peers list

// AddPeer adds a new peer to the list if it's not already present
func AddPeer(peerAddress string) {
	mu.Lock()
	defer mu.Unlock()

	for _, peer := range peers {
		if peer == peerAddress {
			return // Peer is already in the list
		}
	}
	peers = append(peers, peerAddress)
}

// GetPeers returns a copy of the list of connected peers
func GetPeers() []string {
	mu.Lock()
	defer mu.Unlock()

	copyPeers := make([]string, len(peers))
	copy(copyPeers, peers)
	return copyPeers
}

// ========================Message Broadcasting========================

// BroadcastMessage sends a message of the specified type and data to all connected peers
func BroadcastMessage(message Message, senderPeer string) {
	mu.Lock()
	defer mu.Unlock()

	messageJSON, err := SerializeMessage(message)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}

	println(senderPeer)

	for _, peer := range peers {
		if peer == senderPeer {
			continue // Skip the sender peer
		}
		println(peer)
		go func(peer string) {
			address := fmt.Sprintf("%s:8080", peer) //Sending over port 8080 (port 8080 will be used as a transaction message sending port)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Println("Error connecting to peer:", err)
				return
			}
			defer conn.Close()

			fmt.Fprintf(conn, "%s\n", messageJSON)
			println("Propagating message to other peer")
		}(peer)
	}
}

//==========================Block Broadcasting=============================

// BroadcastBlock broadcasts a block to all connected peers except the sender peer
func BroadcastBlock(block *blockchain.Block, senderPeer string) {
	mu.Lock()
	defer mu.Unlock()

	blockJSON, err := json.Marshal(block)
	if err != nil {
		fmt.Println("Error serializing block:", err)
		return
	}

	for _, peer := range peers {
		if peer == senderPeer {
			continue // Skip the sender peer
		}

		go func(peer string) {
			address := fmt.Sprintf("%s:6000", peer) //Sending over port 6000 (port 6000 will be used as a block sending port)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Println("Error connecting to peer:", err)
				return
			}
			defer conn.Close()

			fmt.Fprintf(conn, "%s\n", blockJSON)
		}(peer)
	}
}

// ========================Message Sending========================

// SendDataHashToRandomPeer selects a random peer and sends the provided data hash
func SendDataHashToRandomPeer(message Message) {
	mu.Lock()
	defer mu.Unlock()

	if len(peers) == 0 {
		fmt.Println("No peers available to send data.")
		return
	}

	// Select a random peer
	rand.Seed(time.Now().UnixNano())
	randomPeer := peers[rand.Intn(len(peers))]

	// Send the data hash to the random peer
	go func(peer string) {
		messageJSON, err := SerializeMessage(message)
		if err != nil {
			fmt.Println("Error serializing message:", err)
			return
		}
		address := fmt.Sprintf("%s:%d", peer, 8080) //Add port to use along with address
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println("Error connecting to peer:", err)
			return
		}
		defer conn.Close()

		fmt.Fprintf(conn, "%s\n", messageJSON)
	}(randomPeer)
}
