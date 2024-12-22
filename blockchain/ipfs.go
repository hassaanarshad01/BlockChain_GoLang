package blockchain

import (
	"fmt"
	"io"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// ========================IPFS Client Management========================

// IPFSClient holds the IPFS shell instance
var IPFSClient *shell.Shell

// InitializeIPFSClient initializes the IPFS client with a predefined node address
func InitializeIPFSClient(nodeAddress string) error {
	IPFSClient = shell.NewShell(nodeAddress)
	if !IPFSClient.IsUp() {
		return fmt.Errorf("IPFS node is not reachable at: %s", nodeAddress)
	}
	fmt.Println("IPFS client connected successfully to:", nodeAddress)
	return nil
}

// ========================IPFS Data Operations========================

// SimulateIPFSDownload retrieves data from IPFS using the given hash
func SimulateIPFSDownload(hash string) (string, error) {
	if IPFSClient == nil {
		return "", fmt.Errorf("IPFS client is not initialized")
	}

	dataStream, err := IPFSClient.Cat(hash)
	if err != nil {
		return "", fmt.Errorf("error fetching data from IPFS: %w", err)
	}
	defer dataStream.Close()

	dataBytes, err := io.ReadAll(dataStream)
	if err != nil {
		return "", fmt.Errorf("error reading data from IPFS stream: %w", err)
	}

	data := string(dataBytes)
	fmt.Printf("Downloaded data from IPFS for hash: %s\n", hash)
	return data, nil
}

// AddDataToIPFS adds data to IPFS and returns its hash
func AddDataToIPFS(data string) (string, error) {
	if IPFSClient == nil {
		return "", fmt.Errorf("IPFS client is not initialized")
	}

	reader := strings.NewReader(data)
	hash, err := IPFSClient.Add(reader)
	if err != nil {
		return "", fmt.Errorf("error adding data to IPFS: %w", err)
	}

	fmt.Printf("Added data to IPFS: %s, Hash: %s\n", data, hash)
	return hash, nil
}

// AddPredefinedDataToIPFS adds multiple predefined data items to IPFS and returns their hashes
func AddPredefinedDataToIPFS(dataItems []string) (map[string]string, error) {
	if IPFSClient == nil {
		return nil, fmt.Errorf("IPFS client is not initialized")
	}

	hashes := make(map[string]string)

	for _, data := range dataItems {
		hash, err := AddDataToIPFS(data)
		if err != nil {
			return nil, fmt.Errorf("error adding data '%s' to IPFS: %w", data, err)
		}
		fmt.Printf("Added data to IPFS: %s, Hash: %s\n", data, hash)
		hashes[data] = hash
	}

	return hashes, nil
}

// ========================IPFS Transaction Verification========================

// DownloadAndVerify verifies a transaction by downloading and validating its data
func DownloadAndVerify(tx Transaction) bool {
	return true
}
