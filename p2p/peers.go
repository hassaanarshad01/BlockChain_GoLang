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

//broadcasting
func BroadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	for _, peer := range peers {
		go ConnectToPeer(peer, message)
	}
}
