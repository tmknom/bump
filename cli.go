package bump

import (
	"flag"
	"fmt"
	"io"
)

func Handle(args []string, outStream, errStream io.Writer) error {
	fs := flag.NewFlagSet("bump", flag.ContinueOnError)
	fs.SetOutput(errStream)
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() == 0 {
		return printHelp(errStream)
	}

	switch fs.Arg(0) {
	case "init":
		cmd := &InitCommand{
			args:      args[1:],
			outStream: outStream,
		}
		return cmd.Run()
	case "major":
		cmd := &MajorCommand{
			args:      args[1:],
			outStream: outStream,
		}
		return cmd.Run()
	case "minor":
		cmd := &MinorCommand{
			args:      args[1:],
			outStream: outStream,
		}
		return cmd.Run()
	case "patch":
		cmd := &PatchCommand{
			args:      args[1:],
			outStream: outStream,
		}
		return cmd.Run()
	case "show":
		cmd := &ShowCommand{
			args:      args[1:],
			outStream: outStream,
		}
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
