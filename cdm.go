package main

import (
	"GoTrack/vcs"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gt",
	Short: "GotTrack is a simple Go-based version control system",
}

func Execute() error {
	return rootCmd.Execute()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the version control system",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Mkdir(".gt", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		os.MkdirAll(".gt/objects", 0755)
		fmt.Println(".gt directory and subdirectories created successfully.")
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Save current state",
	Run: func(cmd *cobra.Command, args []string) {
		fileTree := vcs.FileTree()

		fileTree.PrintTree(".")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(commitCmd)
}
