package vcs

import (
	"GoTrack/constants"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func HandleInit(cwd string) {
	gtDir := filepath.Join(cwd, constants.GTDir)
	objectsDir := filepath.Join(cwd, constants.ObjectsDir)
	err := os.Mkdir(gtDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	os.MkdirAll(objectsDir, 0755)
	fmt.Println(".gt directory and subdirectories created successfully.")
}

func HandleCommit(commitMessage string, cwd string) {

	GTDirPath := filepath.Join(cwd, constants.GTDir)
	if _, err := os.Stat(GTDirPath); os.IsNotExist(err) {
		log.Fatal("GoTrack is not initilized.")
		return
	}

	objectsDir := filepath.Join(cwd, constants.ObjectsDir)

	fileTree := RootDir(cwd)

	fileTree.PrintFileTree("")

	tree := BuildTree(fileTree)
	WriteTree(&tree, objectsDir)
	latestCommit, _ := GetLatestCommitHash()

	commit := WriteCommit(tree.Hash, latestCommit, commitMessage, objectsDir)

	UpdateLatestCommit(GTDirPath, commit.Hash)
	UpdateCurrentCommit(GTDirPath, commit.Hash)

}

func HandleLog() {
	latestCommit, err := GetLatestCommitHash()
	if err != nil {
		fmt.Println("Error getting latest commit:", err)
		return
	}

	printCommit(latestCommit)
}

func HandleCheckout(cwd string, hash string) {

	GTDirPath := filepath.Join(cwd, constants.GTDir)

	UpdateCurrentCommit(GTDirPath, hash)
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
