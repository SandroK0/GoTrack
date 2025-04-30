package vcs

import (
	"GoTrack/constants"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Blob struct {
	Hash string
	Name string
}

type Tree struct {
	Hash    string
	Entries []TreeEntry
}

type TreeEntry struct {
	Mode    string
	Type    string // "blob" or "tree"
	Hash    string
	Name    string
	Content []byte
	Entries []TreeEntry // Only for tree
}

func WriteBlob(file *TreeEntry, GTDirPath string) (string, error) {

	blobPath := filepath.Join(GTDirPath, constants.ObjectsDir, file.Hash[:2], file.Hash[2:])

	if _, err := os.Stat(blobPath); err == nil {
		return file.Hash, nil
	}

	if err := os.MkdirAll(filepath.Dir(blobPath), 0755); err != nil {
		return "", err
	}

	outFile, err := os.Create(blobPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err := outFile.Write(file.Content); err != nil {
		return "", err
	}

	return file.Hash, nil
}

func WriteTree(tree *TreeEntry, GTDirPath string) {

	treePath := filepath.Join(GTDirPath, constants.ObjectsDir, tree.Hash[:2], tree.Hash[2:])

	if err := os.MkdirAll(filepath.Dir(treePath), 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(treePath, tree.Content, 0644); err != nil {
		log.Fatal(err)
	}

	for _, entry := range tree.Entries {
		if entry.Type == "tree" {
			WriteTree(&entry, GTDirPath)
		} else {
			WriteBlob(&entry, GTDirPath)
		}
	}

}

func BuildTree(fileTree *Directory) TreeEntry {
	var entries []TreeEntry

	for _, file := range fileTree.Files {
		fileHash := HashContent(file.Content)

		fileContentWithHeader := append([]byte(fmt.Sprintf("blob %d\000", len(file.Content))), file.Content...)
		entries = append(entries, TreeEntry{
			Mode:    "100644", // Regular file
			Type:    "blob",
			Hash:    fileHash,
			Name:    file.Name,
			Content: fileContentWithHeader,
		})
	}

	for _, dir := range fileTree.SubDirs {
		subTree := BuildTree(dir)
		entries = append(entries, TreeEntry{
			Mode: "040000", // Directory mode
			Type: "tree",
			Hash: subTree.Hash,
			Name: dir.Name,
		})
	}
	// fmt.Println("entries:\n", entries)
	for _, entry := range entries {
		fmt.Println("Name:", entry.Name)
	}

	return constructTree(entries)
}

func constructTree(entries []TreeEntry) TreeEntry {
	var treeData []byte

	for _, entry := range entries {

		name := []byte(entry.Name)
		hash := entry.Hash

		treeData = append(treeData, []byte(entry.Mode)...)
		treeData = append(treeData, ' ') // Space separator
		treeData = append(treeData, name...)
		treeData = append(treeData, ' ')
		treeData = append(treeData, hash...)
		treeData = append(treeData, '\n')
	}

	// Create the tree content by adding the header: "tree <size>\0"
	treeContent := append([]byte(fmt.Sprintf("tree %d\000", len(treeData))), treeData...)

	fmt.Println("tree content:\n", string(treeContent))

	treeHash := HashContent(treeData)

	return TreeEntry{Hash: treeHash, Entries: entries, Content: treeContent}
}

func ParseTree(data string, hash string) Tree {

	lines := strings.Split(data, "\n")
	tree := Tree{}

	tree.Hash = hash

	for _, line := range lines {

		treeEntry := TreeEntry{}

		parts := strings.Split(line, " ")
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
