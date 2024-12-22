package blockchain

import (
	"fmt"
)

// ========================Represents a mining node========================
type Miner struct{}

// ========================Mines a block by verifying all transactions and performing proof-of-work========================
func (m *Miner) Mine(transactions []Transaction, prevhash []byte) *Block {

	// Verify each transaction
	for _, tx := range transactions {

		if !DownloadAndVerify(tx) {
			panic(fmt.Sprintf("Transaction verification failed: %+v", tx))
		}

	}

	// Create a new block with verified transactions
	block := NewBlock(transactions, prevhash)

	return block
}
