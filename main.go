package main

import (
	"fmt"
	"os"
)

const (
	GTDir      = ".gt"
	ObjectsDir = ".gt/objects"
	CommitsDir = ".gt/commits" // You can add more paths if necessary
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
