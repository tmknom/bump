package bump

import (
	"flag"
	"fmt"
	"io"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct {
	*baseBumpCommand
}

func newMajorCommand(args []string, outStream, errStream io.Writer) *MajorCommand {
	return &MajorCommand{
		baseBumpCommand: &baseBumpCommand{
			versionType: MAJOR,
			args:        args,
			outStream:   outStream,
			errStream:   errStream,
		},
	}
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct {
	*baseBumpCommand
}

func newMinorCommand(args []string, outStream, errStream io.Writer) *MinorCommand {
	return &MinorCommand{
		baseBumpCommand: &baseBumpCommand{
			versionType: MINOR,
			args:        args,
			outStream:   outStream,
			errStream:   errStream,
		},
	}
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct {
	*baseBumpCommand
}

func newPatchCommand(args []string, outStream, errStream io.Writer) *PatchCommand {
	return &PatchCommand{
		baseBumpCommand: &baseBumpCommand{
			versionType: PATCH,
			args:        args,
			outStream:   outStream,
			errStream:   errStream,
		},
	}
}

type baseBumpCommand struct {
	versionType *VersionType
	args        []string
	outStream   io.Writer
	errStream   io.Writer
}

func (c *baseBumpCommand) run() error {
	fs := flag.NewFlagSet(fmt.Sprintf("bump %s", c.versionType.subcommand()), flag.ContinueOnError)
	fs.SetOutput(c.errStream)

	var versionFile string
	fs.StringVar(&versionFile, "version-file", defaultVersionFile, "A version file for storing current version")

	err := fs.Parse(c.args)
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

	bump := NewBump(version, c.versionType, versionFile, c.outStream)
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

	writer := NewVersionWriter(FileWriteType, version, c.versionFile)
	version, err = writer.Write()
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

	prefix := ""
	fs.StringVar(&prefix, "prefix", prefix, "Show version with prefix")

	err := fs.Parse(c.args)
	if err != nil {
		return err
	}

	file := NewVersionIO(c.versionFile)

	version, err := file.Read()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c.outStream, "%s%s\n", prefix, version.string())
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

	writer := NewVersionWriter(FileWriteType, b.version, b.versionFile)
	return writer.Write()
}

func (b *Bump) upFromCommandLine() (*Version, error) {
	err := b.version.up(b.versionType)
	if err != nil {
		return nil, err
	}

	writer := NewVersionWriter(FileWriteType, b.version, b.versionFile)
	return writer.Write()
}
