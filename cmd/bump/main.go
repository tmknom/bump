package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tmknom/bump"
)

const VersionFile = "VERSION"
const InitialVersion = "0.1.0"

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: bump <subcommand> [flags]")
		os.Exit(0)
	}

	err := handle()
	if err != nil {
		log.Fatalln(err)
	}
}

func handle() error {
	flag.Parse()
	switch flag.Arg(0) {
	case "init":
		cmd := &bump.InitCommand{}
		return cmd.Run(InitialVersion, VersionFile)
	case "show":
		cmd := &bump.ShowCommand{}
		return cmd.Run(VersionFile)
	case "up":
		cmd := &bump.UpCommand{}
		return cmd.Run(VersionFile)
	}
	return nil
}
