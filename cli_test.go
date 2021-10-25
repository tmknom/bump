package bump

import (
	"bytes"
	"os"
	"testing"
)

func TestHandle(t *testing.T) {
	cases := []struct {
		args     []string
		expected string
	}{
		{[]string{"init"}, "0.1.0\n"},
		{[]string{"patch"}, "0.1.1\n"},
		{[]string{"minor"}, "0.2.0\n"},
		{[]string{"major"}, "1.0.0\n"},
	}

	for _, tc := range cases {
		stdout := &bytes.Buffer{}
		err := Handle(tc.args, stdout, os.Stderr)
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", tc.args, err)
		}

		raw, err := os.ReadFile(defaultVersionFile)
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", tc.args, err)
		}

		got := string(raw)
		if got != tc.expected {
			t.Fatalf("%q - unexpected version: got=%q, expected=%q", tc.args, got, tc.expected)
		}

		if stdout.String() != tc.expected {
			t.Fatalf("%q - unexpected stdout: got=%q, expected=%q", tc.args, stdout.String(), tc.expected)
		}
	}
}

func TestHandleSubcommand(t *testing.T) {
	cases := []struct {
		subcommand string
		args       []string
		expected   string
	}{
		{"init", []string{}, "0.1.0\n"},
		{"patch", []string{}, "0.1.1\n"},
		{"minor", []string{}, "0.2.0\n"},
		{"major", []string{}, "1.0.0\n"},
	}

	for _, tc := range cases {
		stdout := &bytes.Buffer{}
		err := handleSubcommand(tc.subcommand, tc.args, stdout)
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", tc.args, err)
		}

		raw, err := os.ReadFile(defaultVersionFile)
		if err != nil {
			t.Fatalf("%q - unexpected error: %s", tc.args, err)
		}

		got := string(raw)
		if got != tc.expected {
			t.Fatalf("%q - unexpected version: got=%q, expected=%q", tc.args, got, tc.expected)
		}

		if stdout.String() != tc.expected {
			t.Fatalf("%q - unexpected stdout: got=%q, expected=%q", tc.args, stdout.String(), tc.expected)
		}
	}
}

func TestHandleSubcommandWithVersion(t *testing.T) {
	cases := []struct {
		subcommand string
		args       []string
		expected   string
	}{
		{"init", []string{"1.2.3"}, "1.2.3\n"},
		{"patch", []string{"1.2.3"}, "1.2.4\n"},
		{"minor", []string{"1.2.3"}, "1.3.0\n"},
		{"major", []string{"1.2.3"}, "2.0.0\n"},
	}

	for _, tc := range cases {
		stdout := &bytes.Buffer{}
		err := handleSubcommand(tc.subcommand, tc.args, stdout)
		if err != nil {
			t.Fatalf("'%q %q' - unexpected error: %s", tc.subcommand, tc.args[0], err)
		}

		if stdout.String() != tc.expected {
			t.Fatalf("'%q %q' - unexpected stdout: got=%q, expected=%q",
				tc.subcommand, tc.args[0], stdout.String(), tc.expected)
		}
	}
}
