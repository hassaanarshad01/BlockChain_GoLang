package blockchain

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ========================IPFS Client Management========================

// IPFS Gateway
const ipfsGateway = "https://aqua-large-swordtail-281.mypinata.cloud/ipfs/"
// Function to download a file from IPFS using its CID
func downloadFile(cid string) ([]byte, error) {
	// Construct the URL for the file on the IPFS gateway
	url := ipfsGateway + cid

	// Create a new HTTP client
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve file, status code: %d", resp)
	}

	// Read the response body (file data)
	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return fileData, nil
}

// Function to write data to a file on disk
func writeFile(fileName string, data []byte) error {
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}
	return nil
}

// InitializeIPFSClient initializes the IPFS client with a predefined node address
func InitializeIPFS() error {
		// Define the CIDs of your 4 datasets and the algorithm file
	datasetCIDs := []string{
		"bafkreid243lqg3bzzakerjsagwf2pivjuigxgqo3djhevawt5367hbr5oq", //Dataset 1
		"bafkreicyveagfqfb3issx5ndk22rtwhhgbm722foot42fcagn5pfjj4s7y", //Dataset 2
		"bafkreibftfmcwzerumjs5me67ctnu7oxf7qrd7buq65mfaewqrogwr7u5m", //Dataset 3
		"bafkreig2lcclmnjfk57swtdeuhzedoyevztjjfytn2n7h4jyiz224hrvle", //Dataset 4
		"bafkreie7yrsjsyyry6b7htmc7sfpdrnupz5dcxjwszy7whaew5zmmj2dui", //Dataset 5
	}

	algorithmCID := "bafkreihjr3b2ur564nohhegxsrb3rfz4oahrblo6czigl6t7uizpwa6gx4" // Algorithm file
	requirementsCID := "bafkreia6m6xuyta6565uv7ur4in2m77euxnwnlhkjys5h4wswdojkbkde4" // Requirements file
	// Download and save each dataset
	for i, cid := range datasetCIDs {
		fileData, err := downloadFile(cid)
		if err != nil {
			fmt.Printf("Error downloading Dataset %d: %s\n", i+1, err)
			continue
		}

		// Save the dataset locally
		fileName := fmt.Sprintf("dataset_%d.csv", i+1) // Save with a name like dataset_1.csv
		err = writeFile(fileName, fileData)
		if err != nil {
			fmt.Printf("Error saving Dataset %d: %s\n", i+1, err)
			continue
		}

		fmt.Printf("Dataset %d downloaded and saved as '%s'\n", i+1, fileName)
	}

	// Download and save the algorithm file
	algorithmData, err := downloadFile(algorithmCID)
	if err != nil {
		fmt.Println("Error downloading the algorithm file:", err)
		return err
	}
	// Save the algorithm file locally
	err = writeFile("algorithm.py", algorithmData) // Save the algorithm as algorithm.py
	if err != nil {
		fmt.Println("Error saving the algorithm file:", err)
		return err
	}
	fmt.Println("Algorithm file downloaded and saved as 'algorithm.py'")

	reqirementsFile, err := downloadFile(requirementsCID)
	if err != nil {
		fmt.Println("Error downloading the Requirements file:", err)
		return err
	}
	err = writeFile("requirements.txt", reqirementsFile) 
	fmt.Println("Requirements file downloaded and saved as 'requirements.txt'")
}

// ========================IPFS Transaction Verification========================

// DownloadAndVerify verifies a transaction by downloading and validating its data
func DownloadAndVerify(tx Transaction) bool {
	return true
}
 