package vcs

import (
	"GoTrack/constants"
	"fmt"
	"os"
)

func HandleInit() {

	err := os.Mkdir(constants.GTDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	os.MkdirAll(constants.ObjectsDir, 0755)
	fmt.Println(".gt directory and subdirectories created successfully.")
}

func HandleCommit(fileTree *Directory, commitMessage string) {

	tree := BuildTree(fileTree)
	fmt.Println("tree name test:", tree.Entries)
	WriteTree(&tree)
	latestCommit, _ := GetLatestCommitHash()

	commit := WriteCommit(tree.Hash, latestCommit, commitMessage)

	UpdateLatestCommit(commit.Hash)
	UpdateCurrentCommit(commit.Hash)

}

func HandleLog() {
	latestCommit, err := GetLatestCommitHash()
	if err != nil {
		fmt.Println("Error getting latest commit:", err)
		return
	}

	printCommit(latestCommit)
}

func HandleCheckout(hash string, fileTree *Directory) {

	UpdateCurrentCommit(hash)
	cleanDirectory(".")
	commitData, _ := ReadObject(hash)
	commit := ParseCommit(string(commitData))
	treeData, _ := ReadObject(commit.TreeHash)
	tree := ParseTree(string(treeData), commit.TreeHash)
	ApplyTree(&tree, ".")
}

func HandleCat(hash string) {
	data, err := ReadObject(hash)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Object Content:")
	fmt.Print(string(data))
}
