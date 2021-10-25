package bump

import (
	"flag"
	"fmt"
	"io"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct {
	args      []string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *MajorCommand) Run() error {
	return parseAndUp("major", c.args, c.outStream, MAJOR)
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct {
	args      []string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *MinorCommand) Run() error {
	return parseAndUp("minor", c.args, c.outStream, MINOR)
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct {
	args      []string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *PatchCommand) Run() error {
	return parseAndUp("patch", c.args, c.outStream, PATCH)
}

func parseAndUp(subcommand string, args []string, outStream io.Writer, versionType VersionType) error {
	fs := flag.NewFlagSet(fmt.Sprintf("bump %s", subcommand), flag.ContinueOnError)
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

	bump := NewBump(version, versionType, outStream)
	return bump.Up()
}

// InitCommand is a command which inits a new version file.
type InitCommand struct {
	args      []string
	outStream io.Writer
}

const defaultInitialVersion = "0.1.0"

// Run runs the procedure of this command.
func (c *InitCommand) Run() error {
	fs := flag.NewFlagSet("bump init", flag.ContinueOnError)
	err := fs.Parse(c.args)
	if err != nil {
		return err
	}

	strVersion := defaultInitialVersion
	if fs.NArg() > 0 {
		strVersion = fs.Arg(0)
	}

	v, err := toVersion(strVersion)
	if err != nil {
		return err
	}

	file := NewVersionIO()
	v, err = file.Write(v)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.outStream, v.string())
	return nil
}

// ShowCommand is a command which show current version.
type ShowCommand struct {
	args      []string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *ShowCommand) Run() error {
	file := NewVersionIO()

	version, err := file.Read()
	if err != nil {
		return err
	}

	fmt.Fprintln(c.outStream, version.string())
	return nil
}

// Bump wraps the basic bump up method.
type Bump struct {
	version     *Version
	versionType VersionType
	outStream   io.Writer
}

// NewBump constructs a new Bump.
func NewBump(version *Version, versionType VersionType, outStream io.Writer) *Bump {
	return &Bump{
		version:     version,
		versionType: versionType,
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
	file := NewVersionIO()

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

	file := NewVersionIO()
	return file.Write(b.version)
}
