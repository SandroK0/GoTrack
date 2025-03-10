package vcs

import (
	"GoTrack/constants"
	"bytes"
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

func ReadCurrent(hash string) ([]byte, error) {

	// Construct the full path to the object file
	objectPath := filepath.Join(constants.CurrentStateDir, hash[:2], hash[2:]) // Store objects in subdirectories like Git

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
