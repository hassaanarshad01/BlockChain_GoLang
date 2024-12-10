package p2p

import (
	"sync"
)

var peers []string
var mu sync.Mutex

//adding a peer if not added before
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

//returns copy of list of peers connected
func GetPeers() []string {
	mu.Lock()
	defer mu.Unlock()

	copyPeers := make([]string, len(peers))
	copy(copyPeers, peers)
	return copyPeers
}

//broadcasting code changed into JSON communication
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

			fmt.Fprintf(conn, messageJSON+"\n")
		}(peer)
	}
}

