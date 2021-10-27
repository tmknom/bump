package bump

import (
	"bytes"
	"os"
	"testing"
)

func TestBaseBumpCommand(t *testing.T) {
	cases := []struct {
		versionType    *VersionType
		args           []string
		expectedFile   string
		expectedStdout string
	}{
		{MINOR, []string{"0.1.0"}, "0.2.0\n", "0.2.0\n"},
		{MINOR, []string{}, "0.3.0\n", "0.3.0\n"},
		{MINOR, []string{"--dry-run"}, "0.3.0\n", "0.4.0\n"},
		{MINOR, []string{"0.1.0", "--dry-run"}, "0.3.0\n", "0.2.0\n"},
	}

	for _, tc := range cases {
		stdout := &bytes.Buffer{}
		command := newBaseBumpCommand(tc.versionType, tc.args, stdout, os.Stderr)

		err := command.run()
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", tc.args, err)
		}

		raw, err := os.ReadFile(defaultVersionFile)
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", tc.args, err)
		}

		got := string(raw)
		if got != tc.expectedFile {
			t.Fatalf("%q - unexpected file version: got=%q, expected=%q", tc.args, got, tc.expectedFile)
		}

		if stdout.String() != tc.expectedStdout {
			t.Fatalf("%q - unexpected stdout version: got=%q, expected=%q", tc.args, stdout.String(), tc.expectedStdout)
		}
	}
}
