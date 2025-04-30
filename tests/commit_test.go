package tests

import (
	"GoTrack/constants"
	"GoTrack/vcs"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type item struct {
	path    string
	isDir   bool
	content string // Only used for files
}

func ValidateFile(t *testing.T, item item, objectsDir string) {
	// Compute hash and expected content for the file
	contentBytes := []byte(item.content)
	fileHash := vcs.HashContent(contentBytes)
	expectedContentWithHeader := append([]byte(fmt.Sprintf("blob %d\000", len(contentBytes))), contentBytes...)

	// Check if file object exists in .gt/objects
	filePath := filepath.Join(objectsDir, fileHash[:2], fileHash[2:])
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file %s object does not exist at %s", item.path, filePath)
	}

	// Verify file object content
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file %s object: %v", item.path, err)
	}
	if string(fileData) != string(expectedContentWithHeader) {
		t.Fatalf("file %s object content mismatch: got %s, want %s", item.path, string(fileData), string(expectedContentWithHeader))
	}
}

func TestCommit(t *testing.T) {
	tmp := t.TempDir()

	// Define an array of items (files and directories)
	items := []item{
		{path: "dir1", isDir: true},
		{path: "dir2", isDir: true},
		{path: "file1.txt", isDir: false, content: "content1"},
		{path: "file2.txt", isDir: false, content: "content2"},
		// {path: "dir2/file3.txt", isDir: false, content: "content3"},
	}

	// Create directories and files
	for _, item := range items {
		fullPath := filepath.Join(tmp, item.path)
		if item.isDir {
			if err := os.Mkdir(fullPath, 0755); err != nil {
				t.Fatalf("failed to create directory %s: %v", item.path, err)
			}
		} else {
			if err := os.WriteFile(fullPath, []byte(item.content), 0644); err != nil {
				t.Fatalf("failed to create file %s: %v", item.path, err)
			}
		}
	}

	// Perform commit
	commitMessage := "test"
	fileTree := vcs.RootDir(tmp)
	vcs.HandleInit(tmp)
	vcs.HandleCommit(fileTree, commitMessage, tmp)

	// Verify .gt/objects directory exists
	gtDir := filepath.Join(tmp, constants.GTDir)
	objectsDir := filepath.Join(gtDir, constants.ObjectsDir)
	fmt.Println("Objects directory:", objectsDir)
	if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
		t.Fatalf(".gt/objects directory does not exist")
	}

	for _, item := range items {
		if item.isDir {
			continue // Skip directories for now (can add tree object checks later if needed)
		}

		// Validate file object
		ValidateFile(t, item, objectsDir)
	}
}
