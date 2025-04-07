package vcs

import (
	"GoTrack/constants"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

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

func ReadStash(hash string) ([]byte, error) {

	// Construct the full path to the object file
	objectPath := filepath.Join(constants.StashDir, hash[:2], hash[2:]) // Store objects in subdirectories like Git

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

func Stash(fileTree *Directory) {
	// Need to write current state in stash
	tree := BuildTree(fileTree)

	fmt.Println("tree entrie:", tree.Hash)

}

func HashContent(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}

func HasUncommitedChanges() {
	// Check for changes (current state vs latest commit)

	// fileTree := RootDir()

}
