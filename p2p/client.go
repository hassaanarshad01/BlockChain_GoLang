package p2p

import (
	"bufio"
	"fmt"
	"net"
)

// ========================Peer-to-Peer Communication========================

// ConnectToPeer establishes a connection with a peer, sends a message, and listens for responses
func ConnectToPeer(peerAddress string, messageType string, dataset interface{}, algo interface{}, req interface{}) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	message := Message{Type: messageType, Dataset: dataset, Algo: algo, Requirements: req}

	// Serialize the message into JSON
	messageJSON, err := SerializeMessage(message)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}

	// Send the serialized message to the peer
	fmt.Fprintf(conn, "%s\n", messageJSON)

	// Listen for responses from the peer
	reader := bufio.NewReader(conn)
	for {
		responseJSON, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by peer:", err)
			break
		}

		// Deserialize the response JSON
		response, err := DeserializeMessage(responseJSON)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			continue
		}

		// Log the response from the peer
		fmt.Println("Response from peer:", response)
	}
}
