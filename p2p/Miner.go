package p2p

import (
	"BlockchainProject/blockchain"
	"BlockchainProject/ipfs"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

var mempool = blockchain.NewMempool()

var ledger, prevHash = blockchain.InitBlockchain()

var peerAddr = GetCurrentMachineAddress()

var stopProcessing chan bool

func init() {
	stopProcessing = make(chan bool)
}

func stopAllProcessing() {
	stopProcessing <- true
}

func GetCurrentMachineAddress() string {

	// Get the IP address of the current machine
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var ipAddr string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddr = ipnet.IP.String()
				break
			}
		}
	}

	return ipAddr
}

func Miner() {
	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening for connections:", err)
		return
	}
	defer ln.Close()

	if err != nil {
		fmt.Println("Error getting machine address:", err)
		return
	}

	fmt.Printf("Peer %s listening on port 8080...", peerAddr)

	//Start checking for mining needs
	go startMiningRoutine()
	go listenForIncomingBlocks()

	for {
		// Accept an incoming connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		defer conn.Close()

		// Handle the connection
		go handleConnection(conn)
	}
}

func handleGeneratorMessage(datasetCID, algoCID, requirementsCID string) error {

	result, err := ipfs.InitializeAndProcess(datasetCID, algoCID, requirementsCID)
	if err != nil {
		return fmt.Errorf("Error running algorithm:", err)
	}

	resultHash, err := ipfs.HashOutput(result)
	if err != nil {
		return fmt.Errorf("Error hashing output:", err)
	}

	fmt.Println("RESULT HASH: ", resultHash)

	//Create transaction
	trans := blockchain.Transaction{
		DataHash:     datasetCID,
		AlgoHash:     algoCID,
		Requirements: requirementsCID,
		Output:       resultHash,
	}

	//Add the transaction to the mempool
	mempool.AddTransaction(&trans)

	return nil
}

func listenForIncomingBlocks() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println("Error listening for incoming blocks:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Listening for incoming blocks on port 6000...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming block connection:", err)
			continue
		}
		defer conn.Close()

		// Handle the incoming block
		go handleIncomingBlock(conn)
	}
}

func handleIncomingBlock(conn net.Conn) {
	// Read the incoming block
	blockJSON, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Error reading incoming block:", err)
		return
	}

	// Unmarshal the block
	var block blockchain.Block
	err = json.Unmarshal(blockJSON, &block)
	if err != nil {
		fmt.Println("Error unmarshaling incoming block:", err)
		return
	}

	// Stop all processing
	stopAllProcessing()

	// Verify the block
	verified, err := VerifyBlock(&block)
	if err != nil {
		fmt.Println("Error verifying block:", err)
		return
	}

	if verified {
		fmt.Println("Block verified successfully. Adding block to ledger")
		ledger.AddBlock(block.Transactions)
		return
	}

	fmt.Println("Block verification failed. Block rejected")

	// Restart the listening for more blocks
	listenForIncomingBlocks()

}

func VerifyBlock(block *blockchain.Block) (bool, error) {
	fmt.Println("Verifying block...")

	// Verify each transaction in the block
	for _, tx := range block.Transactions {
		isVerified, err := ipfs.VerifyTransaction(tx.Output, tx.DataHash, tx.AlgoHash, tx.Requirements)
		if err != nil {
			return false, fmt.Errorf("error verifying transaction: %w", err)
		}
		if !isVerified {
			return false, errors.New("transaction verification failed")
		}
	}

	fmt.Println("Block verified successfully.")
	return true, nil
}

func mineBlock(txs []blockchain.Transaction) {

	// Create a new block with the transactions from the mempool (Mining is done within the NewBlock() function)
	block := blockchain.NewBlock(txs, prevHash)

	//Assign new prevHash to be used for the next block
	prevHash = block.Hash

	// Add the block to the ledger
	ledger.AddBlock(block.Transactions)

	fmt.Println("Block mined and added to ledger: Block Hash->", block.Hash)

	// Propagate the block to all other peers
	fmt.Println("Block propagation started")
	go BroadcastBlock(block, GetCurrentMachineAddress())

	startMiningRoutine()
}

func startMiningRoutine() {
	go func() {
		for {
			select {
			case <-stopProcessing:
				return
			default:
				// Check if the mempool has 5 or more transactions
				println("Transactions in Mempool: ", len(mempool.GetTransactions()))
				if len(mempool.GetTransactions()) >= 2 {

					//Dereferencing pointers
					transactions := mempool.GetTransactions()
					var txs []blockchain.Transaction
					for _, tx := range transactions {
						txs = append(txs, *tx)
					}

					// Start forming and mining the block
					go mineBlock(txs)

					// Clear the mempool
					mempool.ClearTransactions()

				}

				// Sleep for 1 second before checking again
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func handleConnection(conn net.Conn) {
	// Read incoming messages
	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-stopProcessing:
			return
		default:
			if !scanner.Scan() {
				return
			}
			messageJSON := scanner.Text()
			var message Message

			message, err := DeserializeMessage(messageJSON)
			if err != nil {
				fmt.Println("Error unmarshaling message:", err)
				continue
			}
			fmt.Println("Received message from peer:", message)

			//BroadcastMessage(message, peerAddr) //==============================================TODO

			//Extract data from message
			datasetCID := message.Dataset.(string)
			algoCID := message.Algo.(string)
			requirementsCID := message.Requirements.(string)

			//Handles the message sent by Generator peer on port 8080
			err = handleGeneratorMessage(datasetCID, algoCID, requirementsCID)
			if err != nil {
				fmt.Println("Error handling generator message:", err)
			}
		}
	}
}
