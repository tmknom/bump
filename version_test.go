package bump

import (
	"testing"
)

func TestVersionUp(t *testing.T) {
	cases := []struct {
		version     *Version
		versionType VersionType
		expected    string
	}{
		{newVersion(major(0), minor(0), patch(0)), PATCH, "0.0.1"},
		{newVersion(major(0), minor(1), patch(9)), PATCH, "0.1.10"},
		{newVersion(major(1), minor(0), patch(10)), PATCH, "1.0.11"},
		{newVersion(major(0), minor(0), patch(0)), MINOR, "0.1.0"},
		{newVersion(major(0), minor(9), patch(2)), MINOR, "0.10.0"},
		{newVersion(major(1), minor(10), patch(0)), MINOR, "1.11.0"},
		{newVersion(major(0), minor(0), patch(0)), MAJOR, "1.0.0"},
		{newVersion(major(9), minor(9), patch(2)), MAJOR, "10.0.0"},
		{newVersion(major(10), minor(10), patch(0)), MAJOR, "11.0.0"},
	}

	for _, tc := range cases {
		before := tc.version.string()
		err := tc.version.up(tc.versionType)
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", before, err)
		}

		got := tc.version.string()
		if got != tc.expected {
			t.Fatalf("%q bumped version is unexpected: got=%q, expected=%q", before, got, tc.expected)
		}
	}
}
