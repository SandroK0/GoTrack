package vcs

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Name    string
	Content []byte
}

type Directory struct {
	Name    string
	Files   []*File
	SubDirs []*Directory
}

func (d *Directory) PrintTree(prefix string) {
	fmt.Println(prefix + d.Name + "/")
	newPrefix := prefix + "  "
	for _, file := range d.Files {
		fmt.Println(newPrefix + file.Name)
	}
	for _, subdir := range d.SubDirs {
		subdir.PrintTree(newPrefix)
	}
}

func (d *Directory) AddFile(name string, content []byte) {
	newFile := &File{Name: name, Content: content}
	d.Files = append(d.Files, newFile)
}

func (d *Directory) AddSubDir(name string) *Directory {
	newDir := &Directory{Name: name}
	d.SubDirs = append(d.SubDirs, newDir)
	return newDir
}

// We get and return entire file tree
func ScanFileTree(d *Directory, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			if entry.Name() == ".gt" {
				continue
			}
			subDir := d.AddSubDir(entry.Name())
			ScanFileTree(subDir, entryPath)
		} else {
			if entry.Name() == "gt" {
				continue
			}
			data, err := os.ReadFile(entryPath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			d.AddFile(entry.Name(), data)
		}
	}
}

func FileTree() *Directory {
	root := &Directory{Name: "root"}
	ScanFileTree(root, ".")
	return root
}
