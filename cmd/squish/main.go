package main

import (
	"fmt"
	"os"
	"squish/internal/cli"
)

var Version = "development"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Squish version %s\n", Version)
		return
	}

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
