package bump

import (
	"flag"
	"fmt"
	"io"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct {
	versionFile string
	args        []string
	outStream   io.Writer
	errStream   io.Writer
}

// Run runs the procedure of this command.
func (c *MajorCommand) Run() error {
	return parseAndUp(MAJOR, c.args, c.outStream, c.errStream)
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct {
	versionFile string
	args        []string
	outStream   io.Writer
	errStream   io.Writer
}

// Run runs the procedure of this command.
func (c *MinorCommand) Run() error {
	return parseAndUp(MINOR, c.args, c.outStream, c.errStream)
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct {
	versionFile string
	args        []string
	outStream   io.Writer
	errStream   io.Writer
}

// Run runs the procedure of this command.
func (c *PatchCommand) Run() error {
	return parseAndUp(PATCH, c.args, c.outStream, c.errStream)
}

func parseAndUp(versionType *VersionType, args []string, outStream, errStream io.Writer) error {
	fs := flag.NewFlagSet(fmt.Sprintf("bump %s", versionType.subcommand()), flag.ContinueOnError)
	fs.SetOutput(errStream)

	var versionFile string
	fs.StringVar(&versionFile, "version-file", defaultVersionFile, "A version file for storing current version")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	var version *Version
	if fs.NArg() > 0 {
		version, err = toVersion(fs.Arg(0))
		if err != nil {
			return err
		}
	}

	bump := NewBump(version, versionType, versionFile, outStream)
	return bump.Up()
}

// InitCommand is a command which inits a new version file.
type InitCommand struct {
	versionFile string
	args        []string
	outStream   io.Writer
	errStream   io.Writer
}

const defaultInitialVersion = "0.1.0"

// Run runs the procedure of this command.
func (c *InitCommand) Run() error {
	fs := flag.NewFlagSet("bump init", flag.ContinueOnError)
	fs.SetOutput(c.errStream)
	fs.StringVar(&c.versionFile, "version-file", defaultVersionFile, "A version file for storing current version")
	err := fs.Parse(c.args)
	if err != nil {
		return err
	}

	strVersion := defaultInitialVersion
	if fs.NArg() > 0 {
		strVersion = fs.Arg(0)
	}

	version, err := toVersion(strVersion)
	if err != nil {
		return err
	}

	file := NewVersionIO(c.versionFile)
	version, err = file.Write(version)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.outStream, version.string())
	return err
}

// ShowCommand is a command which show current version.
type ShowCommand struct {
	versionFile string
	args        []string
	outStream   io.Writer
	errStream   io.Writer
}

// Run runs the procedure of this command.
func (c *ShowCommand) Run() error {
	fs := flag.NewFlagSet("bump show", flag.ContinueOnError)
	fs.SetOutput(c.errStream)
	fs.StringVar(&c.versionFile, "version-file", defaultVersionFile, "A version file for storing current version")
	err := fs.Parse(c.args)
	if err != nil {
		return err
	}

	file := NewVersionIO(c.versionFile)

	version, err := file.Read()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.outStream, version.string())
	return err
}

// Bump wraps the basic bump up method.
type Bump struct {
	version     *Version
	versionType *VersionType
	versionFile string
	outStream   io.Writer
}

// NewBump constructs a new Bump.
func NewBump(version *Version, versionType *VersionType, versionFile string, outStream io.Writer) *Bump {
	return &Bump{
		version:     version,
		versionType: versionType,
		versionFile: versionFile,
		outStream:   outStream,
	}
}

// Up increments the current version.
func (b *Bump) Up() error {
	var err error
	if b.version != nil {
		b.version, err = b.upFromCommandLine()
	} else {
		b.version, err = b.upFromFile()
	}

	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(b.outStream, b.version.string())
	return err
}

func (b *Bump) upFromFile() (*Version, error) {
	file := NewVersionIO(b.versionFile)

	var err error
	b.version, err = file.Read()
	if err != nil {
		return nil, err
	}

	err = b.version.up(b.versionType)
	if err != nil {
		return nil, err
	}

	return file.Write(b.version)
}

func (b *Bump) upFromCommandLine() (*Version, error) {
	err := b.version.up(b.versionType)
	if err != nil {
		return nil, err
	}

	file := NewVersionIO(b.versionFile)
	return file.Write(b.version)
}
