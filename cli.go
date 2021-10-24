package bump

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const VersionFile = "VERSION"
const InitialVersion = "0.1.0"

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
		return cmd.Run(InitialVersion, VersionFile)
	case "major":
		cmd := &MajorCommand{version: argVersion}
		return cmd.Run(VersionFile)
	case "minor":
		cmd := &MinorCommand{version: argVersion}
		return cmd.Run(VersionFile)
	case "patch":
		cmd := &PatchCommand{version: argVersion}
		return cmd.Run(VersionFile)
	case "show":
		cmd := &ShowCommand{}
		return cmd.Run(VersionFile)
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
