package blockchain

import "math/big"

// PoW represents the Proof-of-Work structure
type PoW struct {
	Block  *Block   // The block for which proof-of-work is being calculated
	target *big.Int // The target value for the proof-of-work difficulty
}

// Sets the number of leading zeros required in the hash (higher = more difficult)
const Difficulty = 14

// NewProof initializes a new PoW instance for a given block
func NewProof(b *Block) *PoW {
	// Create a big integer with value 1
	target := big.NewInt(1)
	// Shift left by (256 - Difficulty) bits to set the difficulty target
	target.Lsh(target, uint(256-Difficulty))
	// Create and return a PoW instance
	pow := &PoW{b, target}
	return pow
}

// ToBytes converts an int64 number into a byte array
func ToBytes(num int64) []byte {
	// Create a buffer to hold the bytes
	var buff = new(bytes.Buffer)

	// Write the integer to the buffer in BigEndian format
	err := binary.Write(buff, binary.BigEndian, num)

	//Error handling
	if err != nil {
		log.Fatal(err)
	}

	// Return integer in byte representation
	return buff.Bytes()
}

// Init prepares the data to be hashed with the given nonce
func (pow *PoW) Init(nonce int) []byte {
	// Join the block data, previous hash, nonce, and difficulty into a single byte array
	data := bytes.Join([][]byte{
		pow.Block.Data,          // Block data
		pow.Block.PrevHash,      // Hash of the previous block
		ToBytes(int64(nonce)),   // Nonce converted to bytes
		ToBytes(int64(Difficulty)) // Difficulty converted to bytes
		},
		[]byte{}
	)

	return data
}

// getHash performs the proof-of-work algorithm to find a valid hash
func (pow *PoW) getHash() (int, []byte) {
	var initHash big.Int // Used to hold the hash as a big integer for comparison
	var hash [32]byte    // The final hash result
	var nonce = 0        // Start with nonce 0

	// Loop to find a valid hash
	for nonce < math.MaxInt64 {
		// Prepare the data with the current nonce
		data := pow.Init(nonce)
		// Compute the SHA-256 hash of the data
		hash = sha256.Sum256(data)
		// Print the hash in hexadecimal format for debugging
		fmt.Printf("%x", hash)
		// Convert the hash to a big integer
		initHash.SetBytes(hash[:])

		// Check if the hash is less than the target (valid PoW)
		if initHash.Cmp(pow.target) == -1 {
			break // Found a valid hash, exit the loop
		} else {
			nonce++ // Increment nonce
		}
	}

	fmt.Println()

	// Return the valid nonce and the corresponding hash
	return nonce, hash[:]
}
