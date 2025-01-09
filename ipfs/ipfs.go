package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"time"
)

// IPFS Gateway (Pinata's Public Gateway)
const ipfsGateway = "https://aqua-large-swordtail-281.mypinata.cloud/ipfs/"

// Struct to hold the JSON output from the Python script
// Including dataset and algorithm CIDs
type AlgorithmResult struct {
	Result struct {
		Centroids [][]float64 `json:"centroids"`
		Inertia   float64     `json:"inertia"`
	} `json:"result"`
	Dataset string `json:"dataset"`
	Algorithm string `json:"algorithm"`
}

// Function to download a file from IPFS using its CID
func downloadFile(cid string) ([]byte, error) {
	url := ipfsGateway + cid
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
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

// Function to write data to a file on disk
func writeFile(fileName string, data []byte) error {
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}
	return nil
}

// Function to execute the Python script and capture its output
func runPythonAlgorithm(scriptName string, dataset string) (string, error) {
	cmdArgs := []string{scriptName, dataset}
	cmd := exec.Command("python", cmdArgs...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run Python script: %w, stderr: %s", err, stderr.String())
	}
	return out.String(), nil
}

// Function to install Python requirements
func installRequirements(requirementsFile string) error {
	cmd := exec.Command("pip", "install", "-r", requirementsFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install requirements: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

func hashOutput(output AlgorithmResult) (string, error) {
    // Marshal JSON with sorted keys to ensure consistency
    outputBytes, err := json.Marshal(output)
    if err != nil {
        return "", fmt.Errorf("error marshaling output for hashing: %w", err)
    }
    // Log the JSON being hashed
    fmt.Printf("Hashing Output (Canonical JSON): %s\n", string(outputBytes))
    hash := sha256.Sum256(outputBytes)
    return hex.EncodeToString(hash[:]), nil
}


// Function to initialize and process a random dataset
func initializeAndProcess(datasetCIDs map[string]string, algorithmCID string, requirementsCID string) (AlgorithmResult, error) {
	fmt.Println("Selecting a random dataset...")
	rand.Seed(time.Now().UnixNano())
	datasetNames := make([]string, 0, len(datasetCIDs))
	for name := range datasetCIDs {
		datasetNames = append(datasetNames, name)
	}
	selectedDataset := datasetNames[rand.Intn(len(datasetNames))]
	datasetCID := datasetCIDs[selectedDataset]

	// Download the selected dataset
	fmt.Printf("Downloading dataset: %s (CID: %s)\n", selectedDataset, datasetCID)
	datasetData, err := downloadFile(datasetCID)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error downloading dataset: %w", err)
	}
	datasetFileName := fmt.Sprintf("%s.csv", selectedDataset)
	err = writeFile(datasetFileName, datasetData)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error saving dataset: %w", err)
	}
	fmt.Printf("Dataset saved as '%s'\n", datasetFileName)

	// Download and save algorithm file
	fmt.Printf("Downloading algorithm file (CID: %s)\n", algorithmCID)
	algorithmData, err := downloadFile(algorithmCID)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error downloading algorithm file: %w", err)
	}
	err = writeFile("algorithm.py", algorithmData)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error saving algorithm file: %w", err)
	}

	// Download and save requirements file
	fmt.Printf("Downloading requirements file (CID: %s)\n", requirementsCID)
	requirementsData, err := downloadFile(requirementsCID)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error downloading requirements file: %w", err)
	}
	err = writeFile("requirements.txt", requirementsData)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error saving requirements file: %w", err)
	}

	// Install requirements
	fmt.Println("Installing Python requirements...")
	err = installRequirements("requirements.txt")
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error installing requirements: %w", err)
	}
	fmt.Println("Python requirements installed successfully.")

	// Run the algorithm
	fmt.Println("Running algorithm...")
	output, err := runPythonAlgorithm("algorithm.py", datasetFileName)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error running Python algorithm: %w", err)
	}

	var result AlgorithmResult
	err = json.Unmarshal([]byte(output), &result.Result)
	if err != nil {
		return AlgorithmResult{}, fmt.Errorf("error parsing JSON output: %w", err)
	}

	// Add dataset and algorithm CIDs to the result
	result.Dataset = datasetCID
	result.Algorithm = algorithmCID

	return result, nil
}

func verifyTransaction(hash string, datasetCID, algorithmCID, requirementsCID string) (bool, error) {
    fmt.Println("Verifying transaction...")

    // Download dataset
    fmt.Printf("Downloading dataset (CID: %s)...\n", datasetCID)
    datasetData, err := downloadFile(datasetCID)
    if err != nil {
        return false, fmt.Errorf("error downloading dataset: %w", err)
    }
    err = writeFile("dataset.csv", datasetData)
    if err != nil {
        return false, fmt.Errorf("error saving dataset: %w", err)
    }

    // Download algorithm
    fmt.Printf("Downloading algorithm (CID: %s)...\n", algorithmCID)
    algorithmData, err := downloadFile(algorithmCID)
    if err != nil {
        return false, fmt.Errorf("error downloading algorithm: %w", err)
    }
    err = writeFile("algorithm.py", algorithmData)
    if err != nil {
        return false, fmt.Errorf("error saving algorithm: %w", err)
    }

    // Download requirements
    fmt.Printf("Downloading requirements (CID: %s)...\n", requirementsCID)
    requirementsData, err := downloadFile(requirementsCID)
    if err != nil {
        return false, fmt.Errorf("error downloading requirements: %w", err)
    }
    err = writeFile("requirements.txt", requirementsData)
    if err != nil {
        return false, fmt.Errorf("error saving requirements: %w", err)
    }

    // Install requirements
    fmt.Println("Installing Python requirements...")
    err = installRequirements("requirements.txt")
    if err != nil {
        return false, fmt.Errorf("error installing requirements: %w", err)
    }

    // Run the algorithm
    fmt.Println("Running algorithm...")
    outputString, err := runPythonAlgorithm("algorithm.py", "dataset.csv")
    if err != nil {
        return false, fmt.Errorf("error running Python algorithm: %w", err)
    }

    fmt.Printf("New Output: %s\n", outputString) // Debug log

    // Parse the new output
    var newOutput AlgorithmResult
    err = json.Unmarshal([]byte(outputString), &newOutput.Result)
    if err != nil {
        return false, fmt.Errorf("error parsing JSON output: %w", err)
    }

    // Add dataset and algorithm CIDs to the new result
    newOutput.Dataset = datasetCID
    newOutput.Algorithm = algorithmCID

    // Hash the new output
    hashedNewOutput, err := hashOutput(newOutput)
    if err != nil {
        return false, fmt.Errorf("error hashing new output: %w", err)
    }

    fmt.Printf("New Hash: %s\n", hashedNewOutput) // Debug log

    // Compare hashes
    if hashedNewOutput == hash {
        fmt.Println("Transaction verified successfully.")
        return true, nil
    }

    fmt.Println("Transaction verification failed.")
    return false, nil
}

// Main function
func main() {
	datasetCIDs := map[string]string{
		"dataset_1": "bafkreid243lqg3bzzakerjsagwf2pivjuigxgqo3djhevawt5367hbr5oq",
		"dataset_2": "bafkreicyveagfqfb3issx5ndk22rtwhhgbm722foot42fcagn5pfjj4s7y",
		"dataset_3": "bafkreibftfmcwzerumjs5me67ctnu7oxf7qrd7buq65mfaewqrogwr7u5m",
		"dataset_4": "bafkreig2lcclmnjfk57swtdeuhzedoyevztjjfytn2n7h4jyiz224hrvle",
		"dataset_5": "bafkreie7yrsjsyyry6b7htmc7sfpdrnupz5dcxjwszy7whaew5zmmj2dui",
	}
	algorithmCID := "bafkreihjr3b2ur564nohhegxsrb3rfz4oahrblo6czigl6t7uizpwa6gx4"
	requirementsCID := "bafkreia6m6xuyta6565uv7ur4in2m77euxnwnlhkjys5h4wswdojkbkde4"

	result, err := initializeAndProcess(datasetCIDs, algorithmCID, requirementsCID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Final Algorithm Result:")
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
    // Generate the hash of the output
    hash, err := hashOutput(result)
    if err != nil {
        fmt.Println("Error hashing output:", err)
        return
    }

    fmt.Printf("Generated Hash: %s\n", hash)

    // Verify the transaction using the hash
    verified, err := verifyTransaction(hash, result.Dataset, result.Algorithm, requirementsCID)
    if err != nil {
        fmt.Println("Verification Error:", err)
    } else if verified {
        fmt.Println("Verification Successful.")
    } else {
        fmt.Println("Verification Failed.")
    }
}
