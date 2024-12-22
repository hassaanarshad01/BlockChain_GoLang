package p2p

import (
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
func BroadcastMessage(messageType string, data interface{}) {
	mu.Lock()
	defer mu.Unlock()

	message := Message{Type: messageType, Data: data}
	messageJSON, err := SerializeMessage(message)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}

	for _, peer := range peers {
		go func(peer string) {
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", err)
				return
			}
			defer conn.Close()

			fmt.Fprintf(conn, "%s\n", messageJSON)
		}(peer)
	}
}

// ========================Message Sending========================

// SendDataHashToRandomPeer selects a random peer and sends the provided data hash
func SendDataHashToRandomPeer(dataHash string) {
	mu.Lock()
	defer mu.Unlock()

	if len(peers) == 0 {
		fmt.Println("No peers available to send data hash.")
		return
	}

	// Select a random peer
	rand.Seed(time.Now().UnixNano())
	randomPeer := peers[rand.Intn(len(peers))]

	// Send the data hash to the random peer
	go func(peer string) {
		message := Message{Type: "DATA_HASH", Data: dataHash}
		messageJSON, err := SerializeMessage(message)
		if err != nil {
			fmt.Println("Error serializing message:", err)
			return
		}
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Println("Error connecting to peer:", err)
			return
		}
		defer conn.Close()

		fmt.Fprintf(conn, "%s\n", messageJSON)
	}(randomPeer)
}
