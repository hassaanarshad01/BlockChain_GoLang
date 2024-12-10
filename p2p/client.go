package p2p

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ConnectToPeer(peerAddress string, messageType string, data interface{}) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	// Serialize and send the message
	message := Message{Type: messageType, Data: data}
	messageJSON, err := SerializeMessage(message)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}
	fmt.Fprintf(conn, messageJSON+"\n")

	// Listen for responses from the peer
	reader := bufio.NewReader(conn)
	for {
		responseJSON, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by peer:", err)
			break
		}

		response, err := DeserializeMessage(responseJSON)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			continue
		}
		fmt.Println("Response from peer:", response)
	}
}
