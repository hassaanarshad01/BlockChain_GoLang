package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
	"math/big"
)

// ========================Proof of Work========================
type PoW struct {
	Block  *Block
	target *big.Int
}

// ========================Creates a new proof-of-work instance========================
func NewProof(b *Block) *PoW {
	// Create a target value for the proof-of-work difficulty
	target := big.NewInt(1)
	target.Lsh(target, uint(256-14))

	// Create and return a PoW instance
	pow := &PoW{b, target}
	return pow
}

// ========================Converts an int64 number into a byte array========================
func ToBytes(num int64) []byte {
	// Create a buffer to hold the bytes
	var buff = new(bytes.Buffer)

	// Write the integer to the buffer in BigEndian format
	err := binary.Write(buff, binary.BigEndian, num)

	// Error handling
	if err != nil {
		log.Fatal(err)
	}

	// Return integer in byte representation
	return buff.Bytes()
}

// ========================Prepares the data to be hashed with the given nonce========================
func (pow *PoW) Init(nonce int) []byte {
	// Convert the slice of Transactions to a JSON byte slice
	transactionsBytes, err := json.Marshal(pow.Block.Transactions)
	if err != nil {
		log.Fatal(err) // or handle the error as you see fit
	}

	data := bytes.Join([][]byte{
		transactionsBytes, // Now a []byte
		[]byte(pow.Block.PrevHash),
		ToBytes(int64(nonce)),
		ToBytes(int64(14)),
	}, []byte{})

	return data
}

// ========================Performs the proof-of-work algorithm to find a valid hash========================
func (pow *PoW) GetHash() (int, string) {
	var initHash big.Int // Used to hold the hash as a big integer for comparison
	var hash string      // The final hash result
	var nonce = 0        // Start with nonce 0

	// Loop to find a valid hash
	for nonce < math.MaxInt64 {
		// Prepare the data with the current nonce
		data := pow.Init(nonce)

		// Compute the SHA-256 hash of the data
		hashBytes := sha256.Sum256(data)
		hash = hex.EncodeToString(hashBytes[:])

		// Convert the hash to a big integer
		initHash.SetBytes(hashBytes[:])

		// Check if the hash is less than the target (valid PoW)
		if initHash.Cmp(pow.target) == -1 {
			break // Found a valid hash, exit the loop
		} else {
			nonce++ // Increment nonce
		}
	}

	// Return the valid nonce and the corresponding hash
	return nonce, hash
}
