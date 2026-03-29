package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *cli.State, cmd cli.Command) error {
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

func HandlerRegister(s *cli.State, cmd cli.Command) error {
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

	fmt.Printf("The user: %s has been created\n", createdUser.Name)
	fmt.Println(createdUser)

	return nil
}

func HandlerListUsers(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
			continue
		}
		fmt.Printf("* %s\n", u.Name)
	}
	return nil
}
