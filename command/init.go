package command

import (
	"os"
)

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
