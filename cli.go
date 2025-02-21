package main

import (
	"GoTrack/constants"
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
		err := os.Mkdir(constants.GTDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		os.MkdirAll(constants.ObjectsDir, 0755)
		fmt.Println(".gt directory and subdirectories created successfully.")
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit <message>",
	Short: "Save current state with a commit message",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (the commit message)
	Run: func(cmd *cobra.Command, args []string) {
		fileTree := vcs.FileTree()
		commitMessage := args[0]
		fmt.Println("Commit message:", commitMessage)
		vcs.HandleCommit(fileTree, commitMessage)
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "See commit history",
	Run: func(cmd *cobra.Command, args []string) {
		vcs.LogHistory()
	},
}

var readCmd = &cobra.Command{
	Use:   "cat <hash>",
	Short: "Read object",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (the commit message)
	Run: func(cmd *cobra.Command, args []string) {

		hash := args[0]

		data, err := vcs.ReadObject(hash)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Object Content:")
		fmt.Print(string(data)) // Print as string for debugging
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(readCmd)
}
