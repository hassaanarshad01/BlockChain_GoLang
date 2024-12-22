package p2p

import (
	"fmt"
	"net"
	"time"
)

var peerHealthCheckInterval = 30 * time.Second // Interval to check peers

func CheckPeerHealth() {
	for {
		time.Sleep(peerHealthCheckInterval)

		mu.Lock()
		for i := len(peers) - 1; i >= 0; i-- {
			peer := peers[i]
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Peer unreachable:", peer)
				// Remove dead peer
				peers = append(peers[:i], peers[i+1:]...)
			} else {
				conn.Close() // Connection is healthy
			}
		}
		mu.Unlock()
	}
}
