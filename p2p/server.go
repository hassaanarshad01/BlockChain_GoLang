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
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by peer")
			break
		}
		fmt.Println("Message received:", message)

		if message == "REQUEST_CHAIN\n" {
			conn.Write([]byte("BLOCKCHAIN_DATA\n"))
		}
	}
}
