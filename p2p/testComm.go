package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Testing communication between Docker containers (peers)
func TestPeerCommunication(peerAddress string) error {

	//Sample message to send
	message := map[string]string{"status": "peer_communication_test", "message": "Hello from peer!"}

	messageBytes, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	//Send the message to the specified peer
	resp, err := http.Post(fmt.Sprintf("http://%s/peer-test", peerAddress), "application/json", bytes.NewBuffer(messageBytes))

	if err != nil {
		return fmt.Errorf("failed to send test message to peer: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	fmt.Printf("Peer communication test with %s successful.\n", peerAddress)
	return nil
}

// Generate and send a transaction message
func GenerateTransactionMessage(peerAddress, algoCID, datasetCID string) error {

	//Create the transaction message
	transaction := map[string]string{
		"algoCID":    algoCID,
		"datasetCID": datasetCID,
	}

	transactionBytes, err := json.Marshal(transaction)

	if err != nil {
		return fmt.Errorf("failed to marshal transaction message: %w", err)
	}

	//============Make it send to any random peer after testing===========--------------------------------------TBD
	//Send the transaction message to the specified peer
	resp, err := http.Post(fmt.Sprintf("http://%s/transaction", peerAddress), "application/json", bytes.NewBuffer(transactionBytes))

	if err != nil {
		return fmt.Errorf("failed to send transaction message to peer: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	fmt.Printf("Transaction message sent to peer %s successfully.\n", peerAddress)
	return nil
}
