package main

import (
	"BlockchainProject/ipfs"
	"BlockchainProject/p2p"
	"fmt"
	"os"
)

func TestIPFS() {

	// Mapping of datasets to their IPFS CIDs
	datasetCIDs := map[string]string{
		"dataset_1": "bafybeicvy4d3odjmys7shzw4cddp6hw2zuoifnbzkiwn6vdomr3wysjkee",
		"dataset_2": "bafkreieqc5e3pzaksbsxo573bkeerfry6lm7qbkg44vsrt3sbix5254koy",
	}

	algorithmCID := "bafkreib22cejsgdjwahqmnpdqfa5hpp6xrpxukeqcqdyn4liazg3i7noku"
	requirementsCID := "bafkreigep2i5w2hw5ek4ufubj4ypdzvyzdhzco3wfpj4b25tjwlclgg5bu"

	// Step 1: Process the datasets and algorithm
	fmt.Println("Starting processing...")
	dataset := p2p.SelectRandomDatasetCID(datasetCIDs)
	result, err := ipfs.InitializeAndProcess(dataset, algorithmCID, requirementsCID)
	if err != nil {
		fmt.Printf("Error during initialization and processing: %v\n", err)
		return
	}

	// Step 2: Generate the hash for the output
	fmt.Println("Hashing the result...")
	hash, err := ipfs.HashOutput(result)
	if err != nil {
		fmt.Printf("Error hashing the output: %v\n", err)
		return
	}

	// Step 3: Verify the transaction
	fmt.Println("Starting verification...")
	isVerified, err := ipfs.VerifyTransaction(hash, result.Dataset, result.Algorithm, requirementsCID)
	if err != nil {
		fmt.Printf("Error verifying the transaction: %v\n", err)
		return
	}

	if isVerified {
		fmt.Println("Transaction verified successfully.")
	} else {
		fmt.Println("Transaction verification failed.")
	}
}

func TestComms() {

	// Add some peers to the list
	p2p.AddPeer("localhost")
	//p2p.AddPeer("192.168.1.101")
	//p2p.AddPeer("192.168.1.102")
	//p2p.AddPeer("192.168.1.103")

	go p2p.Miner()

	p2p.InitMessage()
}

// Main function to demonstrate the process
func main() {

	//Testing of IPFS related functionalities
	//TestIPFS()

	//Testing Communication
	//TestComms()

	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "Gen/MINER")
		return
	}

	role := os.Args[1]
	switch role {
	case "Gen":
		p2p.InitMessage()
	case "MINER":
		go p2p.Miner()
	default:
		fmt.Println("Invalid role. Please use Gen or MINER.")
	}

}
