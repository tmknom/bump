package main

import (
	"flag"
	"log"
	"os"

	"github.com/tmknom/bump"
)

func main() {
	err := bump.Run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil && err != flag.ErrHelp {
		log.Fatalln(err)
	}
}
