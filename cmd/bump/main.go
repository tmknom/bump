package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tmknom/bump"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: bump <subcommand> [flags]")
		os.Exit(0)
	}

	err := bump.Handle()
	if err != nil {
		log.Fatalln(err)
	}
}
