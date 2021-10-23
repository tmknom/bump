package command

import (
	"fmt"
	"os"
)

// ShowCommand is a command which show current version.
type ShowCommand struct{}

// Run runs the procedure of this command.
func (c *ShowCommand) Run(filename string) error {
	version, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, string(version))
	return nil
}
