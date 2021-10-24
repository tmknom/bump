package bump

import (
	"fmt"
	"io"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct {
	version   string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *MajorCommand) Run() error {
	return runBump(c.outStream, c.version, MAJOR)
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct {
	version   string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *MinorCommand) Run() error {
	return runBump(c.outStream, c.version, MINOR)
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct {
	version   string
	outStream io.Writer
}

// Run runs the procedure of this command.
func (c *PatchCommand) Run() error {
	return runBump(c.outStream, c.version, PATCH)
}

// InitCommand is a command which inits a new version file.
type InitCommand struct {
	version   string
	outStream io.Writer
}

const defaultInitialVersion = "0.1.0"

// Run runs the procedure of this command.
func (c *InitCommand) Run() error {
	file := NewVersionIO()

	if len(c.version) == 0 {
		c.version = defaultInitialVersion
	}
	v, err := toVersion(c.version)
	if err != nil {
		return err
	}

	v, err = file.Write(v)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.outStream, v.string())
	return nil
}

// ShowCommand is a command which show current version.
type ShowCommand struct {
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

func runBump(outStream io.Writer, currentVersion string, versionType VersionType) error {
	bump := NewBump(currentVersion, versionType)
	version, err := bump.Up()
	if err != nil {
		return err
	}
	fmt.Fprintln(outStream, version.string())
	return nil
}

// Bump wraps the basic bump up method.
type Bump struct {
	current     string
	versionType VersionType
}

// NewBump constructs a new Bump.
func NewBump(current string, versionType VersionType) *Bump {
	return &Bump{
		current:     current,
		versionType: versionType,
	}
}

// Up increments the current version.
func (b *Bump) Up() (*Version, error) {
	if len(b.current) != 0 {
		return b.upFromCommandLine()
	}
	return b.upFromFile()
}

func (b *Bump) upFromFile() (*Version, error) {
	file := NewVersionIO()

	version, err := file.Read()
	if err != nil {
		return nil, err
	}

	err = version.up(b.versionType)
	if err != nil {
		return nil, err
	}

	return file.Write(version)
}

func (b *Bump) upFromCommandLine() (*Version, error) {
	file := NewVersionIO()

	version, err := toVersion(b.current)
	if err != nil {
		return nil, err
	}

	err = version.up(b.versionType)
	if err != nil {
		return nil, err
	}

	return file.Write(version)
}
