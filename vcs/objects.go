package vcs

import (
	"GoTrack/constants"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TreeEntry represents an entry in a tree (file or directory)
type TreeEntry struct {
	Mode string // "100644" for files, "040000" for directories
	Type string // "blob" or "tree"
	Hash string // SHA-1 hash of the object
	Name string // File or directory name
}

// Tree represents a directory structure
type Tree struct {
	Hash    string
	Entries []TreeEntry
}

// Commit represents a commit object
type Commit struct {
	TreeHash   string
	ParentHash string
	TimeStamp  int64
	Message    string // Commit message
	Hash       string // Commit hash

}

// We need:
// Blob objects for files
// objects for trees
// objects for commits

func WriteFileBlob(file *File, objectsDir string) (string, error) {
	// Prepare the content for hashing by adding the Git-like header
	header := fmt.Sprintf("blob %d\000", len(file.Content))
	contentWithHeader := append([]byte(header), file.Content...)

	// Create the hash of the content including the header
	hasher := sha1.New()
	if _, err := hasher.Write(contentWithHeader); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Define the path where the object will be stored
	blobPath := filepath.Join(objectsDir, hash[:2], hash[2:])

	// Check if the blob already exists
	if _, err := os.Stat(blobPath); err == nil {
		fmt.Println("Blob already exists:", hash)
		return hash, nil
	}

	// Create the directory structure for the object
	if err := os.MkdirAll(filepath.Dir(blobPath), 0755); err != nil {
		return "", err
	}

	// Create the blob file and write the content with header
	outFile, err := os.Create(blobPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Write the content to the file
	if _, err := outFile.Write(contentWithHeader); err != nil {
		return "", err
	}

	// Optionally, print the hash of the blob
	fmt.Printf("Blob created for file '%s' with hash: %s\n", file.Name, hash)
	return hash, nil
}

func HashContent(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}

// WriteTree creates a tree object from directory entries
func WriteTree(entries []TreeEntry, objectsDir string) Tree {
	var treeData []byte

	// Process each tree entry and append it to treeData
	for _, entry := range entries {
		// Convert mode to octal string, ensuring it is properly encoded
		mode := fmt.Sprintf("%o", entry.Mode) // Mode as octal string
		name := []byte(entry.Name)            // Entry name as a byte slice
		hash := entry.Hash                    // The hash is already a byte slice

		// Prepare the entry's binary format: <mode> <name>\0<hash>
		// Append the mode, name, null byte, and hash
		treeData = append(treeData, []byte(mode)...)
		treeData = append(treeData, ' ') // Space separator
		treeData = append(treeData, name...)
		treeData = append(treeData, 0) // Null byte separator
		treeData = append(treeData, hash...)
	}

	// Create the tree content by adding the header: "tree <size>\0"
	treeContent := append([]byte(fmt.Sprintf("tree %d\000", len(treeData))), treeData...)

	// Compute the hash of the tree content
	treeHash := HashContent(treeContent)

	// Create the full object path based on the hash
	treePath := filepath.Join(objectsDir, treeHash[:2], treeHash[2:])

	// Create the necessary directories for the object path
	if err := os.MkdirAll(filepath.Dir(treePath), 0755); err != nil {
		log.Fatal(err) // Handle error appropriately in your code
	}

	// Write the tree content to the object file in binary format
	if err := os.WriteFile(treePath, treeContent, 0644); err != nil {
		log.Fatal(err) // Handle error appropriately in your code
	}

	// Return the tree object with the computed hash and entries
	return Tree{Hash: treeHash, Entries: entries}
}

// BuildTree recursively builds a tree object from a directory
func BuildTree(fileTree *Directory, objectsDir string) Tree {
	var entries []TreeEntry

	// Process all files in the directory
	for _, file := range fileTree.Files {
		fileHash, _ := WriteFileBlob(file, objectsDir)
		entries = append(entries, TreeEntry{
			Mode: "100644", // Regular file
			Type: "blob",   // File content as a blob
			Hash: fileHash,
			Name: file.Name,
		})
	}

	// Recursively process all subdirectories
	for _, dir := range fileTree.SubDirs {
		subTree := BuildTree(dir, objectsDir)
		entries = append(entries, TreeEntry{
			Mode: "040000", // Directory mode
			Type: "tree",   // Subdirectory as a tree
			Hash: subTree.Hash,
			Name: dir.Name,
		})
	}

	// Write the tree object to disk in binary format and return it
	return WriteTree(entries, objectsDir)
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

func UpdateLatestCommit(commitHash string) error {
	headPath := filepath.Join(constants.GTDir, "HEAD")
	return os.WriteFile(headPath, []byte(commitHash+"\n"), 0644)
}

func ReadObject(hash string) ([]byte, error) {
	// Construct the full path to the object file
	objectPath := filepath.Join(constants.ObjectsDir, hash[:2], hash[2:]) // Store objects in subdirectories like Git

	// Read the binary data
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, err
	}

	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid object format: missing header separator")
	}

	// Return the content after the null byte
	return data[nullIndex+1:], nil
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
		return // Stop recursion if no parent
	}

	commitData, err := ReadObject(commitHash)
	if err != nil {
		fmt.Println("Error reading commit:", err)
		return
	}

	commitString := string(commitData)
	commit := ParseCommit(commitString)

	fmt.Printf("Tree: %s\nParent: %s\nTimestamp: %d\nMessage: %s\n",
		commit.TreeHash, commit.ParentHash, commit.TimeStamp, commit.Message)

	fmt.Println("\n------------------------------------------------------")

	// Recursively print parent commits
	printCommit(commit.ParentHash)
}

// Main function to log commit history
func LogHistory() {
	latestCommit, err := GetLatestCommitHash()
	if err != nil {
		fmt.Println("Error getting latest commit:", err)
		return
	}

	printCommit(latestCommit)
}

func HandleCommit(fileTree *Directory, commitMessage string) {

	tree := BuildTree(fileTree, constants.ObjectsDir)

	latestCommit, _ := GetLatestCommitHash()

	commit := WriteCommit(tree.Hash, latestCommit, commitMessage)

	UpdateLatestCommit(commit.Hash)

}
