package bump

import "flag"

const VersionFile = "VERSION"
const InitialVersion = "0.1.0"

func Handle() error {
	flag.Parse()

	var argVersion string
	if flag.NArg() > 1 {
		argVersion = flag.Arg(1)
	}

	switch flag.Arg(0) {
	case "init":
		cmd := &InitCommand{}
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
