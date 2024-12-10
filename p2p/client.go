package p2p

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ConnectToPeer(peerAddress string, message string) {
	//TCP connection
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}

	fmt.Fprintf(conn, message+"\n")

	//goroutine to listen to messages from peer
	go func() {
		reader := bufio.NewReader(conn)
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed by peer:", err)
				break
			}
			fmt.Println("Response from peer:", response)
		}
	}()

	//connection stays open, not closed as it is inefficient otherwise
	for {
		var input string
		fmt.Print("Enter message to send to peer: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}
		fmt.Fprintf(conn, input+"\n")
	}
}
