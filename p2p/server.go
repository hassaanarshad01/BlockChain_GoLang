package p2p

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//TCP server p2p communication
func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on port", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

//incoming messages handle
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

		//Handling the message based on it's type, including JSON
		switch message.Type {
		case "REQUEST_CHAIN":
			//Placeholder blockchain data
			response := Message{Type: "BLOCKCHAIN_DATA", Data: "Genesis Block + First Block"}
			responseJSON, _ := SerializeMessage(response)
			conn.Write([]byte(responseJSON + "\n"))
		case "NEW_BLOCK":
			fmt.Println("New block received:", message.Data)
			//process new block
		default:
			fmt.Println("Unknown message type:", message.Type)
		}
	}
}


