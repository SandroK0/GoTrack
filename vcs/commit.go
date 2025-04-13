package vcs

import (
	"GoTrack/constants"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Commit represents a commit object
type Commit struct {
	TreeHash   string
	ParentHash string
	TimeStamp  int64
	Message    string // Commit message
	Hash       string // Commit hash

}

func WriteCommit(treeHash, parentHash, message string) Commit {
	timestamp := time.Now().Unix()

	// Construct the commit content in binary format
	var commitData []byte

	// Add tree hash
	commitData = append(commitData, []byte(fmt.Sprintf("tree %s\n", treeHash))...)

	// Add parent hash (if there's a parent)
	if parentHash != "" {
		commitData = append(commitData, []byte(fmt.Sprintf("parent %s\n", parentHash))...)
	}

	// Add timestamp
	commitData = append(commitData, []byte(fmt.Sprintf("timestamp %d\n", timestamp))...)

	// Add commit message (ensure the message is properly encoded in binary)
	commitData = append(commitData, []byte(fmt.Sprintf("message %s\n", message))...)

	// Create the final commit content by including the header: "commit <size>\0"
	commitContent := append([]byte(fmt.Sprintf("commit %d\000", len(commitData))), commitData...)

	// Compute the hash of the commit content
	commitHash := HashContent(commitContent)

	// Create the full object path based on the hash
	commitPath := filepath.Join(constants.ObjectsDir, commitHash[:2], commitHash[2:])

	// Create the necessary directories for the object path
	if err := os.MkdirAll(filepath.Dir(commitPath), 0755); err != nil {
		log.Fatal(err) // Handle error appropriately in your code
	}

	if err := os.WriteFile(commitPath, commitContent, 0644); err != nil {
		// Write the commit content to the object file in binary format
		log.Fatal(err) // Handle error appropriately in your code
	}

	// Return the commit object with the computed hash and content
	return Commit{
		TreeHash:   treeHash,
		ParentHash: parentHash,
		TimeStamp:  timestamp,
		Message:    message,
		Hash:       commitHash,
	}
}

func GetLatestCommitHash() (string, error) {
	headPath := filepath.Join(constants.GTDir, "HEAD")

	data, err := os.ReadFile(headPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No commits yet
		}
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func GetCurrentCommitHash() (string, error) {
	headPath := filepath.Join(constants.GTDir, "CURRENT")

	data, err := os.ReadFile(headPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No commits yet
		}
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func UpdateCurrentCommit(commitHash string) error {
	headPath := filepath.Join(constants.GTDir, "HEAD")
	return os.WriteFile(headPath, []byte(commitHash+"\n"), 0644)
}

func UpdateLatestCommit(commitHash string) error {
	headPath := filepath.Join(constants.GTDir, "HEAD")
	return os.WriteFile(headPath, []byte(commitHash+"\n"), 0644)
}

func ParseCommit(data string) Commit {
	lines := strings.Split(data, "\n") // Split into lines
	commit := Commit{}

	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2) // Split each line into key-value pair
		if len(parts) < 2 {
			continue // Skip empty or malformed lines
		}

		key, value := parts[0], parts[1]

		switch key {
		case "tree":
			commit.TreeHash = value
		case "parent":
			commit.ParentHash = value
		case "timestamp":
			timestamp, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				commit.TimeStamp = timestamp
			}
		case "message":
			commit.Message = value
		}
	}

	return commit
}

// Recursive function to print commit history
func printCommit(commitHash string) {
	if commitHash == "" {
		return
	}

	commitData, err := ReadObject(commitHash)
	if err != nil {
		fmt.Println("Error reading commit:", err)
		return
	}

	commitString := string(commitData)
	commit := ParseCommit(commitString)

	fmt.Printf("\nHash: %s\nTree: %s\nParent: %s\nTimestamp: %d\nMessage: %s\n",
		commitHash, commit.TreeHash, commit.ParentHash, commit.TimeStamp, commit.Message)
	fmt.Println("\n------------------------------------------------------")

	printCommit(commit.ParentHash)
}
