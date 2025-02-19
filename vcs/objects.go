package vcs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
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
	TreeHash    string
	ParentHash  string
	CreatedTime time.Time
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
	var treeData []string
	for _, entry := range entries {
		line := fmt.Sprintf("%s %s %s %s", entry.Mode, entry.Type, entry.Hash, entry.Name)
		treeData = append(treeData, line)
	}

	treeContent := strings.Join(treeData, "\n")
	treeHash := HashContent([]byte(treeContent))

	tree := Tree{Hash: treeHash, Entries: entries}
	_ = os.WriteFile(objectsDir+treeHash, []byte(treeContent), 0644)

	return tree
}

// BuildTree recursively builds a tree object from a directory
func BuildTree(fileTree *Directory, objectsDir string) Tree {
	var entries []TreeEntry

	for _, file := range fileTree.Files {
		fileHash, _ := WriteFileBlob(file, objectsDir)
		entries = append(entries, TreeEntry{
			Mode: "100644",
			Type: "blob",
			Hash: fileHash,
			Name: file.Name,
		})
	}

	for _, dir := range fileTree.SubDirs {
		subTree := BuildTree(dir, objectsDir)
		entries = append(entries, TreeEntry{
			Mode: "040000",
			Type: "tree",
			Hash: subTree.Hash,
			Name: dir.Name,
		})
	}
	return WriteTree(entries, objectsDir)
}

func HandleCommit(fileTree *Directory, objectsDir string) {

	BuildTree(fileTree, objectsDir)
}
