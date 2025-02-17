package main

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Function to create the .gt directory
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

// Function to hash and compress a file for version control
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Save curent state",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		fmt.Printf("Adding file: %s\n", file)
		// Hash the file
		hash, err := hashFile(file)
		if err != nil {
			fmt.Println("Error hashing file:", err)
			return
		}
		// Compress the file (or the hashed data)
		err = compressFile(file, hash)
		if err != nil {
			fmt.Println("Error compressing file:", err)
			return
		}
		fmt.Printf("File %s added with hash: %x\n", file, hash)
	},
}

// Function to hash the file content (SHA-256)
func hashFile(file string) ([]byte, error) {
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a SHA-256 hash object
	hash := sha256.New()
	// Copy the file contents into the hash function
	_, err = io.Copy(hash, f)
	if err != nil {
		return nil, err
	}
	// Return the file's hash
	return hash.Sum(nil), nil
}

// Function to compress the file using GZIP
func compressFile(file string, hash []byte) error {
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a new GZIP file to store the compressed file
	compressedFileName := fmt.Sprintf(".gt/objects/%x.gz", hash) // Using the hash as filename
	outFile, err := os.Create(compressedFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a GZIP writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Copy the file contents to the GZIP writer
	_, err = io.Copy(gzipWriter, f)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	// Add commands to the root command
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(initCmd)
}

var rootCmd = &cobra.Command{
	Use:   "gt",
	Short: "GotTrack is a simple Go-based version control system",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
