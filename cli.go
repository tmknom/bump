package bump

import (
	"flag"
	"fmt"
	"io"
)

func Run(args []string, outStream, errStream io.Writer) error {
	fs := flag.NewFlagSet("bump", flag.ContinueOnError)
	fs.SetOutput(errStream)
	version := fs.Bool("version", false, "Show version")
	fs.Usage = func() { printTopLevelUsage(fs.Output()) }

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if *version {
		return printVersion(errStream)
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return flag.ErrHelp
	}

	return runSubcommand(fs.Arg(0), args[1:], outStream, errStream)
}

func runSubcommand(subcommand string, args []string, outStream, errStream io.Writer) error {
	switch subcommand {
	case "init":
		cmd := newInitCommand(args, outStream, errStream)
		return cmd.Run()
	case "major":
		cmd := newMajorCommand(args, outStream, errStream)
		return cmd.run()
	case "minor":
		cmd := newMinorCommand(args, outStream, errStream)
		return cmd.run()
	case "patch":
		cmd := newPatchCommand(args, outStream, errStream)
		return cmd.run()
	case "show":
		cmd := newShowCommand(args, outStream, errStream)
		return cmd.Run()
	}
	return nil
}

func printTopLevelUsage(out io.Writer) {
	message := `Bump version that following semantic versioning

Usage:
  bump <command> [<version>] [flags]

Commands:
  init              Init version
  major             Bump up major version
  minor             Bump up minor version
  patch             Bump up patch version
  show              Show the current version

Flags:
  --help            Show help for command
  --version         Show bump version

Examples:
  $ bump init
  $ bump patch
  $ bump minor 1.0.0
`
	printUsage(out, message)
}

func printUsage(out io.Writer, message string) {
	// skip error handling because of caller method is never handled error
	// see detail: usage method inside flag package
	_, _ = fmt.Fprint(out, message)
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "bump v%s\n", "0.0.1")
	return err
}
