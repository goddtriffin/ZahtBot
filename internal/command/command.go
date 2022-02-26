package command

import "fmt"

const commandSignal = "/"

// Command is a Discord bot command.
type Command struct {
	Name        string
	Description string
}

func (c Command) String() string {
	return fmt.Sprintf("%s%s", commandSignal, c.Name)
}

// Help returns the name of the command followed by its description.
func (c Command) Help() string {
	return fmt.Sprintf("%s: %s", c.String(), c.Description)
}
