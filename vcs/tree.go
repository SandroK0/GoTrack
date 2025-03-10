package vcs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func HashContent(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}

// WriteTree creates a tree object from directory entries
func WriteTree(entries []TreeEntry, objectsDir string) Tree {
	var treeData []byte

	// Process each tree entry and append it to treeData
	for _, entry := range entries {

		name := []byte(entry.Name) // Entry name as a byte slice
		hash := entry.Hash         // The hash is already a byte slice

		// Prepare the entry's binary format: <mode> <name>\0<hash>
		// Append the mode, name, null byte, and hash
		treeData = append(treeData, []byte(entry.Mode)...)
		treeData = append(treeData, ' ') // Space separator
		treeData = append(treeData, name...)
		treeData = append(treeData, ' ')
		treeData = append(treeData, hash...)

		treeData = append(treeData, '\n')
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

func ParseTree(data string, hash string) Tree {

	lines := strings.Split(data, "\n") // Split into lines

	tree := Tree{}

	tree.Hash = hash

	for _, line := range lines {

		treeEntry := TreeEntry{}

		parts := strings.Split(line, " ") // Split each line into key-value pair
		if len(parts) < 3 {
			continue
		}
		switch parts[0] {
		case "100644":
			treeEntry.Mode = "100644"
			treeEntry.Type = "blob"
		case "040000":
			treeEntry.Mode = "040000"
			treeEntry.Type = "tree"
		}
		treeEntry.Name = parts[1]
		treeEntry.Hash = parts[2]
		tree.Entries = append(tree.Entries, treeEntry)
	}

	return tree
}

func PrintTree(tree Tree) {

	fmt.Println("Tree Hash:", tree.Hash)

	for _, entry := range tree.Entries {

		fmt.Println("Name:", entry.Name)
		fmt.Println("Mode:", entry.Mode)
		fmt.Println("Type:", entry.Type)
		fmt.Println("Hash:", entry.Hash)

	}

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
