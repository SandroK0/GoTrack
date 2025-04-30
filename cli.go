package main

import (
	"GoTrack/vcs"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gt",
	Short: "GoTrack is a simple Go-based version control system",
}

func Execute() error {
	return rootCmd.Execute()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the version control system",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory:", err)
			return
		}
		vcs.HandleInit(cwd)
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit <message>",
	Short: "Save current state with a commit message",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (the commit message)
	Run: func(cmd *cobra.Command, args []string) {

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory:", err)
			return
		}
		fileTree := vcs.RootDir(cwd)
		commitMessage := args[0]
		fmt.Println("Commit message:", commitMessage)
		vcs.HandleCommit(fileTree, commitMessage, cwd)
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "See commit history",
	Run: func(cmd *cobra.Command, args []string) {
		vcs.HandleLog()
	},
}

var catCmd = &cobra.Command{
	Use:   "cat <hash>",
	Short: "Read object",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (the commit message)
	Run: func(cmd *cobra.Command, args []string) {

		hash := args[0]

		vcs.HandleCat(hash)
	},
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout <hash>",
	Short: "Checkout",
	Args:  cobra.ExactArgs(1), // Expect exactly one argument (the commit message)
	Run: func(cmd *cobra.Command, args []string) {

		// Should check for uncommited changes here.
		// if True. user should commit or stash
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory:", err)
			return
		}
		hash := args[0]
		fileTree := vcs.RootDir(cwd)
		vcs.HandleCheckout(hash, fileTree)
	},
}

var stashCmd = &cobra.Command{
	Use:   "stash",
	Short: "Stash",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var stashApplyCmd = &cobra.Command{
	Use:   "stash apply",
	Short: "Back to current uncommited state",
	Run: func(cmd *cobra.Command, args []string) {

		tree, _ := vcs.ReadStash(args[0])
		tree_struct := vcs.ParseTree(string(tree), args[0])
		vcs.PrintTree(tree_struct)
	},
}

var testCmd = &cobra.Command{
	Use:   "test <hash>",
	Short: "See commit history",
	Run: func(cmd *cobra.Command, args []string) {
		vcs.HasUncommitedChanges()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(catCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(stashCmd)
	rootCmd.AddCommand(stashApplyCmd)
	rootCmd.AddCommand(testCmd)
}
