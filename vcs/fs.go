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

func (d *Directory) PrintDir(prefix string) {
	fmt.Println(prefix + d.Name + "/")
	newPrefix := prefix + "  "
	for _, file := range d.Files {
		fmt.Println(newPrefix + file.Name)
	}
	for _, subdir := range d.SubDirs {
		subdir.PrintDir(newPrefix)
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
func ScanDir(d *Directory, path string) {
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
			ScanDir(subDir, entryPath)
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

func CreateFile(file *File, path string) error {

	// Write the content to the file
	if err := os.WriteFile(path, file.Content, 0644); err != nil {
		return err
	}

	return nil
}

func ApplyTree(tree *Tree, path string) {

	for _, entry := range tree.Entries {
		fullPath := filepath.Join(path, entry.Name)

		switch entry.Type {
		case "blob":
			fileContent, err := ReadObject(entry.Hash)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			file := File{entry.Name, fileContent}
			CreateFile(&file, fullPath)

		case "tree":
			os.Mkdir(fullPath, os.ModePerm)
			treeData, _ := ReadObject(entry.Hash)

			subTree := ParseTree(string(treeData), entry.Hash)
			ApplyTree(&subTree, fullPath)

		}

	}

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

func RootDir(path string) *Directory {
	root := &Directory{Name: "root"}
	ScanDir(root, path)
	return root
}
