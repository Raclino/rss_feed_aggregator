package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/config"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmd map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Cmd[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmd[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: login <username>")
	}

	username := cmd.Args[0]
	ctx := context.Background()

	_, err := s.Db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("user %q does not exist", username)
	}

	if err := s.Config.SetUser(username); err != nil {
		return err
	}

	fmt.Println("The user has been set")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: register <username>")
	}

	username := cmd.Args[0]
	ctx := context.Background()

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	createdUser, err := s.Db.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}

	if err := s.Config.SetUser(createdUser.Name); err != nil {
		return err
	}

	fmt.Println("The user has been created")
	fmt.Println(createdUser)

	return nil
}
