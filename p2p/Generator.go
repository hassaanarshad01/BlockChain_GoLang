package p2p

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// Function to initialize and send a message every 10 seconds
func InitMessage() {
	// Mapping of datasets to their IPFS CIDs
	datasetCIDs := map[string]string{
		"dataset_1": "bafybeicvy4d3odjmys7shzw4cddp6hw2zuoifnbzkiwn6vdomr3wysjkee",
		"dataset_2": "bafkreieqc5e3pzaksbsxo573bkeerfry6lm7qbkg44vsrt3sbix5254koy",
	}

	algorithmCID := "bafkreib22cejsgdjwahqmnpdqfa5hpp6xrpxukeqcqdyn4liazg3i7noku"
	requirementsCID := "bafkreigep2i5w2hw5ek4ufubj4ypdzvyzdhzco3wfpj4b25tjwlclgg5bu"

	// Infinite loop to send messages every 10 seconds
	for {
		// Select a random dataset CID
		randomDatasetCID := SelectRandomDatasetCID(datasetCIDs)

		// Construct the message
		message := Message{
			Type:         "TRANS",
			Dataset:      randomDatasetCID,
			Algo:         algorithmCID,
			Requirements: requirementsCID,
		}

		SendDataHashToRandomPeer(message)

		// Wait for 10 seconds before generating the next message
		time.Sleep(10 * time.Second)
	}
}

// Function to randomly select a dataset CID from the map
func SelectRandomDatasetCID(datasetCIDs map[string]string) string {
	datasets := make([]string, 0, len(datasetCIDs))
	for _, cid := range datasetCIDs {
		datasets = append(datasets, cid)
	}
	return datasets[rand.Intn(len(datasets))]
}

// Function to send the serialized message to a peer
func SendToPeer(peerAddress, message string) {

	address := fmt.Sprintf("%s:%d", peerAddress, 8080) //Add port to use along with address
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	// Send the message
	fmt.Fprintf(conn, "%s\n", message)
	fmt.Println("Message sent to peer:", message)
}
