package main

import (
	"flag"
	"log"
	"os"

	"github.com/tmknom/bump"
)

func main() {
	err := bump.Handle(os.Args[1:], os.Stdout)
	if err != nil && err != flag.ErrHelp {
		log.Fatalln(err)
	}
}
