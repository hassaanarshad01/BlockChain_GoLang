package p2p

import (
	"net"
	"sync"
)

var peers []string
var mu sync.Mutex

//adding new peers to the list
func AddPeer(peerAddress string) {
	mu.Lock()
	defer mu.Unlock()
	peers = append(peers, peerAddress)
}

//broadcasting system
func BroadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	for _, peer := range peers {
		go ConnectToPeer(peer, message)
	}
}
