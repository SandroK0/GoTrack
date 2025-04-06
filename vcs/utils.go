package vcs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

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

func cleanDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == dir {
			return nil
		}

		if info.Name() == "gt" || (info.IsDir() && info.Name() == ".gt") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}

		return nil
	})
}

func HandleCommit(fileTree *Directory, commitMessage string) {

	tree := BuildTree(fileTree)

	latestCommit, _ := GetLatestCommitHash()

	commit := WriteCommit(tree.Hash, latestCommit, commitMessage)

	UpdateLatestCommit(commit.Hash)
	UpdateCurrentCommit(commit.Hash)

}

func LogHistory() {
	latestCommit, err := GetLatestCommitHash()
	if err != nil {
		fmt.Println("Error getting latest commit:", err)
		return
	}

	printCommit(latestCommit)
}

func Checkout(hash string, fileTree *Directory) {

	UpdateCurrentCommit(hash)
	cleanDirectory(".")
	commitData, _ := ReadObject(hash)
	commit := ParseCommit(string(commitData))
	treeData, _ := ReadObject(commit.TreeHash)
	tree := ParseTree(string(treeData), commit.TreeHash)
	ApplyTree(&tree, ".")
}
