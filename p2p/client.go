package p2p

import (
	"bufio"
	"fmt"
	"net"
)

//Connecting to peer
func ConnectToPeer(peerAddress string, message string) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	//message sending
	fmt.Fprintf(conn, message+"\n")

	//recieving response
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Response from peer:", response)
}
