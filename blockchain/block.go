package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// ========================Represents a transaction in the blockchain========================
type Transaction struct {
	DataHash string // Hash of the dataset stored on IPFS
	AlgoHash string // Hash of the AI algorithm stored on IPFS
	Output   string // Expected output of the algorithm
}

// ========================Represents a block in the blockchain========================
type Block struct {
	Hash         []byte        // Hash of the block
	Transactions []Transaction // Transactions stored in the block
	PrevHash     []byte        // Hash of the previous block
	Nonce        int           // Nonce used in proof-of-work
}

// ========================Blockchain========================
type Blockchain struct {
	Blocks []*Block // Slice of blocks forming the blockchain
}

// ========================Calculates and sets the hash for the block========================
func (b *Block) GetHash() {
	// Serialize transactions into a single byte slice
	txData, _ := json.Marshal(b.Transactions)

	// Join transaction data and previous hash
	info := bytes.Join([][]byte{txData, b.PrevHash}, []byte{})

	// Compute the SHA-256 hash of the joined data
	hash := sha256.Sum256(info)

	// Store the hash in the block
	b.Hash = hash[:]
}

// ========================Creates a new block========================
func NewBlock(transactions []Transaction, prevhash []byte) *Block {
	block := &Block{[]byte{}, transactions, prevhash, 0}
	pow := NewProof(block)
	nonce, hash := pow.getHash()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// ========================Initializes the blockchain with a genesis block========================
func InitBlockchain() *Blockchain {
	// Transaction stored in the Genesis Block
	genesisTx := []Transaction{
		{DataHash: "GenesisDataHash", AlgoHash: "GenesisAlgoHash", Output: "GenesisOutput"},
	}

	return &Blockchain{[]*Block{NewBlock(genesisTx, []byte{})}}
}

// ========================Adds a new block to the blockchain========================
func (chain *Blockchain) AddBlock(transactions []Transaction) {
	prevB := chain.Blocks[len(chain.Blocks)-1]
	newB := NewBlock(transactions, prevB.Hash)
	chain.Blocks = append(chain.Blocks, newB)
}

// ========================Get the latest block========================
func (chain *Blockchain) GetLatestBlock() *Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

// ========================Add a block to the chain========================
func (chain *Blockchain) AddBlockToChain(block *Block) {
	chain.Blocks = append(chain.Blocks, block)
}

// ========================Simple deterministic algorithm========================
func DeterministicAlgorithm(data string) string {
	var result byte
	for _, b := range []byte(data) {
		result ^= b
	}
	return hex.EncodeToString([]byte{result})
}

// ========================Hash data using SHA-256========================
func HashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
