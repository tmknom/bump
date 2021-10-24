package bump

import (
	"fmt"
	"os"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct {
	version string
}

// Run runs the procedure of this command.
func (c *MajorCommand) Run(filename string) error {
	return runBump(c.version, filename, MAJOR)
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct {
	version string
}

// Run runs the procedure of this command.
func (c *MinorCommand) Run(filename string) error {
	return runBump(c.version, filename, MINOR)
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct {
	version string
}

// Run runs the procedure of this command.
func (c *PatchCommand) Run(filename string) error {
	return runBump(c.version, filename, PATCH)
}

// InitCommand is a command which inits a new version file.
type InitCommand struct{}

// Run runs the procedure of this command.
func (c *InitCommand) Run(version string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(version + "\n")
	if err != nil {
		return err
	}

	return nil
}

func runBump(currentVersion string, filename string, versionType VersionType) error {
	bump := NewBump(currentVersion, filename, versionType)
	version, err := bump.Up()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, version.string())
	return nil
}

// Bump wraps the basic bump up method.
type Bump struct {
	current     string
	path        string
	versionType VersionType
}

// NewBump constructs a new Bump.
func NewBump(current string, path string, versionType VersionType) *Bump {
	return &Bump{
		current:     current,
		path:        path,
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
	file := NewVersionIO(b.path)

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
	file := NewVersionIO(b.path)

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
