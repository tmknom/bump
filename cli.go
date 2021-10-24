package bump

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func Handle(argv []string) error {
	if len(argv) == 0 {
		return printHelp(os.Stderr)
	}

	flag.Parse()

	var argVersion string
	if flag.NArg() > 1 {
		argVersion = flag.Arg(1)
	}

	switch flag.Arg(0) {
	case "init":
		cmd := &InitCommand{version: argVersion}
		return cmd.Run()
	case "major":
		cmd := &MajorCommand{version: argVersion}
		return cmd.Run()
	case "minor":
		cmd := &MinorCommand{version: argVersion}
		return cmd.Run()
	case "patch":
		cmd := &PatchCommand{version: argVersion}
		return cmd.Run()
	case "show":
		cmd := &ShowCommand{}
		return cmd.Run()
	}
	return nil
}

func printHelp(out io.Writer) error {
	_, err := fmt.Fprintln(out, "Usage: bump <subcommand> [<version>] [flags]")
	if err != nil {
		return err
	}
	return flag.ErrHelp
}
