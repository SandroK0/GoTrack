package main

import (
	_ "compress/gzip"
	_ "crypto/sha256"
	"fmt"
	_ "io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the version control system",
	Run: func(cmd *cobra.Command, args []string) {
		// Create .gt directory to store VCS files
		err := os.Mkdir(".gt", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		// Create subdirectories for storing objects, etc.
		os.MkdirAll(".gt/objects", 0755)
		fmt.Println(".gt directory and subdirectories created successfully.")
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Save curent state",
	Run: func(cmd *cobra.Command, args []string) {
		fileTree := fileTree()

		fmt.Print(fileTree)

	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(initCmd)
}

var rootCmd = &cobra.Command{
	Use:   "gt",
	Short: "GotTrack is a simple Go-based version control system",
}

type File struct {
	Name    string // Name of the file
	Content []byte // File contents
}

// Directory represents a directory which can contain files and subdirectories.
type Directory struct {
	Name    string       // Name of the directory
	Files   []*File      // Files in this directory
	SubDirs []*Directory // Subdirectories in this directory
}

func (d *Directory) printTree(prefix string) {
	fmt.Println(prefix + d.Name + "/")
	newPrefix := prefix + "  "
	for _, file := range d.Files {
		fmt.Println(newPrefix + file.Name)
	}
	for _, subdir := range d.SubDirs {
		subdir.printTree(newPrefix)
	}
}

// AddFile dynamically adds a new file to the directory.
func (d *Directory) AddFile(name string, content []byte) {
	newFile := &File{Name: name, Content: content}
	d.Files = append(d.Files, newFile)
}

// AddSubDir dynamically adds a new subdirectory to the directory.
func (d *Directory) AddSubDir(name string) *Directory {
	newDir := &Directory{Name: name}
	d.SubDirs = append(d.SubDirs, newDir)
	return newDir
}

func buildTree(d *Directory, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			subDir := d.AddSubDir(entry.Name())
			buildTree(subDir, entryPath)
		} else {
			data, err := os.ReadFile(entryPath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			d.AddFile(entry.Name(), data)
		}
	}
}

func fileTree() *Directory {
	root := &Directory{Name: "project"}
	buildTree(root, ".") // Start from current directory
	root.printTree(".")

	return root
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
