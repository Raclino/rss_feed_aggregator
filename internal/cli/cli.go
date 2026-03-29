package cli

import (
	"fmt"
)

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}
