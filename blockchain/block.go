package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// ========================Represents a transaction in the blockchain========================
type Transaction struct {
	DataHash     string // CID of the dataset stored on IPFS
	AlgoHash     string // CID of the AI algorithm stored on IPFS
	Requirements string // CID of the requirements file stored on IPFS
	Output       string // Hash of expected output of the algorithm
}

// ========================Represents a block in the blockchain========================
type Block struct {
	Hash         string        // Hash of the block
	Transactions []Transaction // Transactions stored in the block
	PrevHash     string        // Hash of the previous block
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
	info := bytes.Join([][]byte{txData, []byte(b.PrevHash)}, []byte{})

	// Compute the SHA-256 hash of the joined data
	hash := sha256.Sum256(info)

	// Store the hash in the block
	b.Hash = hex.EncodeToString(hash[:])
}

// ========================Creates a new block========================
func NewBlock(transactions []Transaction, prevhash string) *Block {
	block := &Block{"", transactions, prevhash, 0}
	pow := NewProof(block)
	nonce, hash := pow.GetHash()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

// ========================Initializes the blockchain with a genesis block========================
func InitBlockchain() (*Blockchain, string) {
	// Transaction stored in the Genesis Block
	genesisTx := []Transaction{
		{DataHash: "GenesisData", AlgoHash: "GenesisAlgo", Requirements: "GenesisReq", Output: "GenesisOutputHash"},
	}

	block := NewBlock(genesisTx, "")

	return &Blockchain{[]*Block{block}}, block.Hash
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

// ========================Hash data using SHA-256========================
func HashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
