package tests

import (
	"GoTrack/constants"
	"GoTrack/vcs"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type Mode string

const (
	File      Mode = "100644"
	Directory Mode = "040000"
)

type item struct {
	path    string
	isDir   bool
	content string
	entries []item
}

func ValidateFile(t *testing.T, item item, objectsDir string) {
	// Compute hash and expected content for the file
	contentBytes := []byte(item.content)
	fileHash := vcs.HashContent(contentBytes)
	expectedContentWithHeader := append([]byte(fmt.Sprintf("blob %d\000", len(contentBytes))), contentBytes...)

	// Check if file object exists in .gt/objects
	fmt.Println("Validating file:", objectsDir, fileHash)
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

func createStructure(t *testing.T, basePath string, item item) {
	fullPath := filepath.Join(basePath, item.path)
	if item.isDir {
		if err := os.Mkdir(fullPath, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", item.path, err)
		}
		for _, child := range item.entries {
			createStructure(t, fullPath, child)
		}
	} else {
		if err := os.WriteFile(fullPath, []byte(item.content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", item.path, err)
		}
	}
}

func ValidateTree(t *testing.T, item item, objectsDir string) string {
	// For each entry, get its hash and mode
	var treeContent []byte
	for _, child := range item.entries {
		var mode, hash string
		if child.isDir {
			mode = string(Directory)
			hash = ValidateTree(t, child, objectsDir)
		} else {
			mode = string(File)
			contentBytes := []byte(child.content)
			hash = vcs.HashContent(contentBytes)
			ValidateFile(t, child, objectsDir)
		}
		// Format: mode SP name NUL hash (as raw bytes)
		entry := fmt.Sprintf("%s %s\x00", mode, filepath.Base(child.path))
		entryBytes := append([]byte(entry), decodeHex(hash)...) // hash as raw bytes
		treeContent = append(treeContent, entryBytes...)
	}
	// Tree object: header + content
	header := fmt.Sprintf("tree %d\x00", len(treeContent))
	treeObj := append([]byte(header), treeContent...)
	treeHash := vcs.HashContent(treeContent)
	// Check if tree object exists
	treePath := filepath.Join(objectsDir, treeHash[:2], treeHash[2:])
	if _, err := os.Stat(treePath); os.IsNotExist(err) {
		t.Fatalf("tree %s object does not exist at %s", item.path, treePath)
	}
	// Verify tree object content
	fileData, err := os.ReadFile(treePath)
	if err != nil {
		t.Fatalf("failed to read tree %s object: %v", item.path, err)
	}
	if string(fileData) != string(treeObj) {
		t.Fatalf("tree %s object content mismatch", item.path)
	}
	return treeHash
}

func decodeHex(s string) []byte {
	b := make([]byte, len(s)/2)
	for i := 0; i < len(b); i++ {
		fmt.Sscanf(s[2*i:2*i+2], "%02x", &b[i])
	}
	return b
}

func TestCommit(t *testing.T) {
	tmp := t.TempDir()

	// Define a nested structure of items (files and directories)
	items := []item{
		{
			path: "dir1", isDir: true, entries: []item{
				{path: "fileA.txt", isDir: false, content: "A"},
				{path: "subdir", isDir: true, entries: []item{
					{path: "fileB.txt", isDir: false, content: "B"},
				}},
			},
		},
		{path: "file1.txt", isDir: false, content: "content1"},
		{path: "file2.txt", isDir: false, content: "content2"},
	}

	// Recursively create directories and files
	for _, it := range items {
		createStructure(t, tmp, it)
	}

	// Perform commit
	commitMessage := "test"

	vcs.HandleInit(tmp)
	vcs.HandleCommit(commitMessage, tmp)

	// Verify .gt/objects directory exists
	objectsDir := filepath.Join(tmp, constants.ObjectsDir)
	if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
		t.Fatalf(".gt/objects directory does not exist")
	}

	// Recursively validate all objects
	for _, it := range items {
		if it.isDir {
			ValidateTree(t, it, objectsDir)
		} else {
			ValidateFile(t, it, objectsDir)
		}
	}
}
