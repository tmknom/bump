package bump

import (
	"fmt"
	"os"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct{}

// Run runs the procedure of this command.
func (c *MajorCommand) Run(filename string) error {
	return runBump(filename, MAJOR)
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct{}

// Run runs the procedure of this command.
func (c *MinorCommand) Run(filename string) error {
	return runBump(filename, MINOR)
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct{}

// Run runs the procedure of this command.
func (c *PatchCommand) Run(filename string) error {
	return runBump(filename, PATCH)
}

func runBump(filename string, versionType VersionType) error {
	bump := NewBump(filename, versionType)
	version, err := bump.Up()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, version.string())
	return nil
}

// Bump wraps the basic bump up method.
type Bump struct {
	path        string
	versionType VersionType
}

// NewBump constructs a new Bump.
func NewBump(path string, versionType VersionType) *Bump {
	return &Bump{
		path:        path,
		versionType: versionType,
	}
}

// Up increments the current version.
func (b *Bump) Up() (*Version, error) {
	file := NewVersionFile(b.path)

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
