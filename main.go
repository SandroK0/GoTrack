package main

import (
	"fmt"
	_ "fmt"
	"os"
	_ "os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
