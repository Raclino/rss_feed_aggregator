package cli

import (
	"context"
	"fmt"
	"os"
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
		return fmt.Errorf("error, %s is not registered", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	_, ok := c.Cmd[name]
	if !ok {
		c.Cmd[name] = f
	}
	fmt.Printf("%s command has been registered !\n", name)
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error, no commands argument provided")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("The user has been set")
	return nil

}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error, no commands argument provided")
	}

	ctx := context.Background()
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[1],
	}
	createdUsr, err := s.Db.CreateUser(ctx, newUser)
	if err != nil {
		defer os.Exit(1)
		return fmt.Errorf("User already exist\n")
	}

	err = s.Config.SetUser(createdUsr.Name)
	if err != nil {
		return err
	}

	fmt.Println("The user has been created!\n")
	fmt.Println(createdUsr)
	return nil

}
