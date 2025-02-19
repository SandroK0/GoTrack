package vcs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// We need:
// Blob objects for files
// objects for trees
// objects for commits

func CreateFileBlob(filePath string, objectsDir string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	blobPath := filepath.Join(objectsDir, hash)

	if _, err := os.Stat(blobPath); err == nil {
		fmt.Println("Blob already exists:", hash)
		return hash, nil
	}

	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return "", err
	}

	file.Seek(0, 0)
	outFile, err := os.Create(blobPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", err
	}

	fmt.Println("Blob created:", hash)
	return hash, nil
}
