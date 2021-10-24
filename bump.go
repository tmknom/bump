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

// VersionFile wraps the I/O method for the version file.
type VersionFile struct {
	path string
}

// NewVersionFile constructs a new VersionFile.
func NewVersionFile(path string) *VersionFile {
	return &VersionFile{
		path: path,
	}
}

// Read reads the version file and returns the current version.
func (f *VersionFile) Read() (*Version, error) {
	bytes, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	return toVersion(string(bytes))
}

// Write writes the version to the version file.
func (f *VersionFile) Write(version *Version) (*Version, error) {
	file, err := os.Create(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.WriteString(version.string() + "\n")
	if err != nil {
		return nil, err
	}
	return version, nil
}

type VersionType int

const (
	_ VersionType = iota
	MAJOR
	MINOR
	PATCH
)

// Version takes the form X.Y.Z: X is the major version, Y is the minor version, and Z is the patch version.
type Version struct {
	major
	minor
	patch
}

type major int
type minor int
type patch int

func toVersion(version string) (*Version, error) {
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

	return newVersion(major(x), minor(y), patch(z)), nil
}

func newVersion(x major, y minor, z patch) *Version {
	return &Version{
		major: x,
		minor: y,
		patch: z,
	}
}

func (v *Version) string() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func (v *Version) up(t VersionType) error {
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

func (v *Version) upMajor() {
	v.major += 1
	v.minor = 0
	v.patch = 0
}

func (v *Version) upMinor() {
	v.minor += 1
	v.patch = 0
}

func (v *Version) upPatch() {
	v.patch += 1
}
