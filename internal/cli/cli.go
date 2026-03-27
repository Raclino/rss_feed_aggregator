package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/Raclino/rss_feed_aggregator/internal/rss"
	"github.com/google/uuid"
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

	fmt.Printf("The user: %s has been created\n", createdUser.Name)
	fmt.Println(createdUser)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	if err := s.Db.DeleteAllUsers(ctx); err != nil {
		return err
	}

	fmt.Println("Database Reset was successful")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
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

func HandlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	url := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(ctx, url)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: <feedName> <feedURL>")
	}
	ctx := context.Background()

	user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return err
	}

	feedName := cmd.Args[0]
	url := cmd.Args[1]
	feedId := uuid.New()
	newFeed := database.CreateFeedParams{
		ID:        feedId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	}
	f, err := s.Db.CreateFeed(ctx, newFeed)
	if err != nil {
		return err
	}

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedId,
	}
	_, err = s.Db.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		return err
	}

	fmt.Println(f)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	ctx := context.Background()

	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Println(feeds)
	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: <feedURL>")
	}

	ctx := context.Background()
	user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := s.Db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		return err
	}
	fmt.Println(feed.Name)

	fmt.Println(user.Name)
	return nil
}

func HandlerFollowing(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: <feedURL>")
	}

	ctx := context.Background()
	user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollowerCurrentUser, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)

	for _, feed := range feedFollowerCurrentUser {
		fmt.Println(feed.Name)
	}

	return nil
}
