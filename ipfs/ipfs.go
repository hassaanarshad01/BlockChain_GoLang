package ipfs

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

// IPFS Gateway for accessing files stored on IPFS
const ipfsGateway = "https://fuchsia-official-porcupine-775.mypinata.cloud/ipfs/"

// Struct to parse the JSON output from the Python script
type AlgorithmResult struct {
	Result struct {
		Centroids [][]float64 `json:"centroids"` // Centroids from clustering algorithm
		Inertia   float64     `json:"inertia"`   // Inertia metric (sum of squared distances)
	} `json:"result"`
	Dataset   string `json:"dataset"`   // CID of the dataset used
	Algorithm string `json:"algorithm"` // CID of the algorithm used
}

// Function to download a file from IPFS using its CID
func DownloadFile(cid string) ([]byte, error) {

	url := ipfsGateway + cid
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	// Add authorization header with API key
	apiKey := "Bearer a0bde220f72a645ce158"
	req.Header.Add("Authorization", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve file, status code: %d", resp.StatusCode)
	}
	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}
	return fileData, nil
}

// Writes downloaded data to a file on disk
func WriteFile(fileName string, data []byte) error {
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}
	return nil
}

// Executes a Python script and captures its output
func RunPythonAlgorithm(scriptName string, dataset string) (string, error) {
	println(scriptName, " ", dataset, "\n")
	//cmdArgs := []string{scriptName, dataset}           // Pass script name and dataset as arguments
	cmd := exec.Command("python", scriptName, dataset) // Run the Python script
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run Python script: %w, stderr: %s", err, stderr.String())
	}

	// println(out.String())

	return out.String(), nil
}

// Installs Python requirements using pip
func InstallRequirements(requirementsFile string) error {
	cmd := exec.Command("pip", "install", "-r", requirementsFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install requirements: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

// Generates a hash from the algorithm result
func HashOutput(output AlgorithmResult) (string, error) {

	// Convert the output to canonical JSON format
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("error marshaling output for hashing: %w", err)
	}
	fmt.Printf("Hashing Output (Canonical JSON): %s\n", string(outputBytes))
	hash := sha256.Sum256(outputBytes)

	return hex.EncodeToString(hash[:]), nil
}

// Downloads a random dataset, algorithm, and requirements from IPFS and processes them
func InitializeAndProcess(datasetCID string, algorithmCID string, requirementsCID string) (AlgorithmResult, error) {

	// Download the dataset
	fmt.Printf("Downloading dataset with CID: %s\n", datasetCID)
	datasetData, err := DownloadFile(datasetCID)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error downloading dataset: %w", err)
	}
	datasetFileName := fmt.Sprintf("%s.csv", "datset")
	err = WriteFile(datasetFileName, datasetData)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error saving dataset: %w", err)
	}
	fmt.Printf("Dataset saved as '%s'\n", datasetFileName)

	// Download and save the algorithm
	fmt.Printf("Downloading algorithm file (CID: %s)\n", algorithmCID)
	algorithmData, err := DownloadFile(algorithmCID)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error downloading algorithm file: %w", err)
	}
	err = WriteFile("algorithm.py", algorithmData)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error saving algorithm file: %w", err)
	}

	// Download and save the requirements file
	fmt.Printf("Downloading requirements file (CID: %s)\n", requirementsCID)
	requirementsData, err := DownloadFile(requirementsCID)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error downloading requirements file: %w", err)
	}
	err = WriteFile("requirements.txt", requirementsData)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error saving requirements file: %w", err)
	}

	// Install Python requirements
	fmt.Println("Installing Python requirements...")
	err = InstallRequirements("requirements.txt")
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error installing requirements: %w", err)
	}
	fmt.Println("Python requirements installed successfully.")

	// Run the algorithm with the dataset
	fmt.Println("Running algorithm...")
	output, err := RunPythonAlgorithm("algorithm.py", datasetFileName)

	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error running Python algorithm: %w", err)
	}

	// Parse the output into the AlgorithmResult struct
	var result AlgorithmResult
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error parsing JSON output: %w", err)
	}

	// Add dataset and algorithm CIDs to the result
	result.Dataset = datasetCID
	result.Algorithm = algorithmCID

	return result, nil
}

// Verifies the output by re-executing the algorithm and comparing results
func VerifyTransaction(hash string, datasetCID, algorithmCID, requirementsCID string) (bool, error) {
	fmt.Println("Verifying transaction...")

	// Download the dataset
	fmt.Printf("Downloading dataset (CID: %s)...\n", datasetCID)
	datasetData, err := DownloadFile(datasetCID)
	if err != nil {
		return false, fmt.Errorf("error downloading dataset: %w", err)
	}
	err = WriteFile("dataset.csv", datasetData)
	if err != nil {
		return false, fmt.Errorf("error saving dataset: %w", err)
	}

	// Download the algorithm
	fmt.Printf("Downloading algorithm (CID: %s)...\n", algorithmCID)
	algorithmData, err := DownloadFile(algorithmCID)
	if err != nil {
		return false, fmt.Errorf("error downloading algorithm: %w", err)
	}
	err = WriteFile("algorithm.py", algorithmData)
	if err != nil {
		return false, fmt.Errorf("error saving algorithm: %w", err)
	}

	// Download the requirements
	fmt.Printf("Downloading requirements (CID: %s)...\n", requirementsCID)
	requirementsData, err := DownloadFile(requirementsCID)
	if err != nil {
		return false, fmt.Errorf("error downloading requirements: %w", err)
	}
	err = WriteFile("requirements.txt", requirementsData)
	if err != nil {
		return false, fmt.Errorf("error saving requirements: %w", err)
	}

	// Install requirements
	fmt.Println("Installing Python requirements...")
	err = InstallRequirements("requirements.txt")
	if err != nil {
		return false, fmt.Errorf("error installing requirements: %w", err)
	}

	// Run the algorithm and compare results
	fmt.Println("Running algorithm...")
	outputString, err := RunPythonAlgorithm("algorithm.py", "dataset.csv")
	if err != nil {
		return false, fmt.Errorf("error running Python algorithm: %w", err)
	}

	// Parse the new output
	var newOutput AlgorithmResult
	err = json.Unmarshal([]byte(outputString), &newOutput)
	if err != nil {
		return false, fmt.Errorf("error parsing JSON output: %w", err)
	}

	// Add dataset and algorithm CIDs to the new result
	newOutput.Dataset = datasetCID
	newOutput.Algorithm = algorithmCID

	// Hash the new output
	hashedNewOutput, err := HashOutput(newOutput)
	if err != nil {
		return false, fmt.Errorf("error hashing new output: %w", err)
	}

	fmt.Printf("=========COMPARE | %s == %s |============", hash, hashedNewOutput)

	// Compare the original hash with the newly generated hash
	if hashedNewOutput == hash {
		fmt.Println("Transaction verified successfully.")
		return true, nil
	}

	fmt.Println("Transaction verification failed.")
	return false, nil
}
