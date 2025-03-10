package vcs

import (
	"GoTrack/constants"
	"fmt"
	"os"
	"path/filepath"
)

func SaveCurrentStateTemp(fileTree *Directory) {

	tree := BuildTree(fileTree, constants.CurrentStateDir)

	fmt.Println("tree entrie:", tree.Hash)

}

func cleanDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dir {
			return nil
		}

		// Skip the "gt" file and ".gt" directory
		if info.Name() == "gt" || (info.IsDir() && info.Name() == ".gt") {
			if info.IsDir() {
				return filepath.SkipDir // Prevents descending into .gt
			}
			return nil
		}

		// Remove file or directory
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}

		return nil
	})
}

func recreate(treeHash string) {

	data, _ := ReadObject(treeHash)

	fmt.Println(string(data))

}

func HandleCommit(fileTree *Directory, commitMessage string) {

	tree := BuildTree(fileTree, constants.ObjectsDir)

	latestCommit, _ := GetLatestCommitHash()

	commit := WriteCommit(tree.Hash, latestCommit, commitMessage)

	UpdateLatestCommit(commit.Hash)

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

	SaveCurrentStateTemp(fileTree)

	cleanDirectory(".")

	recreate(hash)

}
