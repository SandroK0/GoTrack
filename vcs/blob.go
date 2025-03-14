package vcs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

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
	// fmt.Printf("Blob created for file '%s' with hash: %s\n", file.Name, hash)
	return hash, nil
}
