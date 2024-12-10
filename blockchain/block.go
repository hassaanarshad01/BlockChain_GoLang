package blockchain

import "bytes"
import "crypto/sha256"

// ================================NOTE================================
// |  Need to change "Data" to Algo and Dataset Hash storing instead  |
// ================================NOTE================================


// Block structure definition
type Block struct {
	Hash     []byte // The hash of the current block
	Data     []byte // Data stored in the block 			(To change)
	PrevHash []byte // Hash of the previous block
	Nonce    int    // Nonce used in proof-of-work
}

// Blockchain represents the chain of blocks
type Blockchain struct {
	Blocks []*Block // Slice of blocks forming the blockchain
}

// GetHash calculates and sets the hash for the block
func (b *Block) GetHash() {
	// Join block data and previous hash into a single byte slice
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	// Compute the SHA-256 hash of the joined data
	hash := sha256.Sum256(info)
	// Store the hash in the block
	b.Hash = hash[:]
}

// NewBlock creates a new block with the given data and the previous block's hash
func NewBlock(data string, prevhash []byte) *Block {
	// Initialize a new block with the provided data and previous hash
	block := &Block{[]byte{}, []byte(data), prevhash, 0}
	// Create a proof-of-work instance for the new block
	pow := NewProof(block)
	// Perform the proof-of-work to find a valid nonce and hash
	nonce, hash := pow.GetHash()
	// Set the block's hash and nonce
	block.Hash = hash[:]
	block.Nonce = nonce

	// Return the newly created block
	return block
}

// InitBlockchain initializes the blockchain with a genesis block
func InitBlockchain() *Blockchain {
	// Create a new blockchain and start with a "Genesis" block
	return &Blockchain{[]*Block{NewBlock("Genesis", []byte{})}}
}

// AddBlock adds a new block with the provided data to the blockchain
func (chain *Blockchain) AddBlock(data string) {
	// Retrieve the last block in the chain
	prevB := chain.Blocks[len(chain.Blocks)-1]
	// Create a new block with the given data and the previous block's hash
	newB := NewBlock(data, prevB.Hash)
	// Append the new block to the blockchain
	chain.Blocks = append(chain.Blocks, newB)
}
