package bump

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MajorCommand is a command which bump up major version.
type MajorCommand struct{}

// Run runs the procedure of this command.
func (c *MajorCommand) Run(filename string) error {
	cmd := &UpCommand{}
	return cmd.Run(filename, MAJOR)
}

// MinorCommand is a command which bump up minor version.
type MinorCommand struct{}

// Run runs the procedure of this command.
func (c *MinorCommand) Run(filename string) error {
	cmd := &UpCommand{}
	return cmd.Run(filename, MINOR)
}

// PatchCommand is a command which bump up patch version.
type PatchCommand struct{}

// Run runs the procedure of this command.
func (c *PatchCommand) Run(filename string) error {
	cmd := &UpCommand{}
	return cmd.Run(filename, PATCH)
}

// UpCommand is a command which bump up version.
type UpCommand struct{}

// Run runs the procedure of this command.
func (c *UpCommand) Run(filename string, versionType VersionType) error {
	version, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	versioning, err := toVersioning(string(version))
	if err != nil {
		return err
	}

	err = versioning.up(versionType)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(versioning.string() + "\n")
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, versioning.string())
	return nil
}

type VersionType int

const (
	_ VersionType = iota
	MAJOR
	MINOR
	PATCH
)

// Versioning takes the form X.Y.Z: X is the major version, Y is the minor version, and Z is the patch version.
type Versioning struct {
	major
	minor
	patch
}

type major int
type minor int
type patch int

func toVersioning(version string) (*Versioning, error) {
	versions := strings.Split(strings.TrimSpace(version), ".")
	x, err := strconv.Atoi(versions[0])
	if err != nil {
		return nil, err
	}

	y, err := strconv.Atoi(versions[1])
	if err != nil {
		return nil, err
	}

	z, err := strconv.Atoi(versions[2])
	if err != nil {
		return nil, err
	}

	return newVersioning(major(x), minor(y), patch(z)), nil
}

func newVersioning(x major, y minor, z patch) *Versioning {
	return &Versioning{
		major: x,
		minor: y,
		patch: z,
	}
}

func (v *Versioning) string() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func (v *Versioning) up(t VersionType) error {
	switch t {
	case MAJOR:
		v.upMajor()
	case MINOR:
		v.upMinor()
	case PATCH:
		v.upPatch()
	default:
		return fmt.Errorf("invalid VersionType: %d", t)
	}
	return nil
}

func (v *Versioning) upMajor() {
	v.major += 1
	v.minor = 0
	v.patch = 0
}

func (v *Versioning) upMinor() {
	v.minor += 1
	v.patch = 0
}

func (v *Versioning) upPatch() {
	v.patch += 1
}
