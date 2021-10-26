package bump

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultVersionFile = "VERSION"

// VersionIO wraps the I/O method for the version file.
type VersionIO struct {
	path string
}

// NewVersionIO constructs a new VersionIO.
func NewVersionIO(path string) *VersionIO {
	return &VersionIO{
		path: path,
	}
}

// Read reads the version file and returns the current version.
func (f *VersionIO) Read() (*Version, error) {
	bytes, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	return toVersion(string(bytes))
}

// Write writes the version to the version file.
func (f *VersionIO) Write(version *Version) (v *Version, err error) {
	file, err := os.Create(f.path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) { err = f.Close() }(file)

	_, err = file.WriteString(version.string() + "\n")
	if err != nil {
		return nil, err
	}
	return version, nil
}

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

func (v *Version) up(t *VersionType) error {
	switch t.value {
	case MajorVersionType:
		v.upMajor()
	case MinorVersionType:
		v.upMinor()
	case PatchVersionType:
		v.upPatch()
	default:
		return fmt.Errorf("invalid versionType: %q", t)
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

type VersionType struct {
	value string
}

func (t *VersionType) subcommand() string {
	return t.value
}

const (
	MajorVersionType = "major"
	MinorVersionType = "minor"
	PatchVersionType = "patch"
)

var MAJOR = &VersionType{value: MajorVersionType}
var MINOR = &VersionType{value: MinorVersionType}
var PATCH = &VersionType{value: PatchVersionType}
