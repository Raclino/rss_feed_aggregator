package cli

import (
	"fmt"

	"github.com/Raclino/rss_feed_aggregator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	login []string
}
type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	// c.cmd[cmd.login[0]]
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.cmd[name]
	if !ok {
		c.cmd[name] = f
	}
	fmt.Println("%s command has been registered !", name)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.login) == 0 {
		return fmt.Errorf("error, no commands argument provided")
	}

	err := s.config.SetUser(cmd.login[0])
	if err != nil {
		return err
	}

	fmt.Println("The user has been set")
	return nil

}
