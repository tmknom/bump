package bump

import "flag"

const VersionFile = "VERSION"
const InitialVersion = "0.1.0"

func Handle() error {
	flag.Parse()
	switch flag.Arg(0) {
	case "init":
		cmd := &InitCommand{}
		return cmd.Run(InitialVersion, VersionFile)
	case "major":
		cmd := &MajorCommand{}
		return cmd.Run(VersionFile)
	case "minor":
		cmd := &MinorCommand{}
		return cmd.Run(VersionFile)
	case "patch":
		cmd := &PatchCommand{}
		return cmd.Run(VersionFile)
	case "show":
		cmd := &ShowCommand{}
		return cmd.Run(VersionFile)
	}
	return nil
}
