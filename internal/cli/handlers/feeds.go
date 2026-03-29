package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/Raclino/rss_feed_aggregator/internal/services"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	ctx := context.Background()

	feedName := cmd.Args[0]
	url := cmd.Args[1]
	feedID := uuid.New()

	newFeed := database.CreateFeedParams{
		ID:        feedID,
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
		FeedID:    feedID,
	}

	_, err = s.Db.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		return err
	}

	fmt.Println(f)
	return nil
}
func HandlerListFeeds(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()

	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Println(feeds)
	return nil
}

func HandlerFollow(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: follow <feed_url>")
	}

	ctx := context.Background()

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

	createdFollow, err := s.Db.CreateFeedFollow(ctx, newFeedFollow)
	if err != nil {
		return err
	}

	fmt.Println(createdFollow.FeedName)
	fmt.Println(createdFollow.UserName)
	return nil
}

func HandlerUnFollow(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	ctx := context.Background()
	feed, err := s.Db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	reqArgs := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	_, err = s.Db.DeleteFeedFollow(ctx, reqArgs)
	if err != nil {
		return err
	}

	fmt.Printf("correctly unfollowed %s feed\n", feed.Name)
	return nil
}

func HandlerFollowing(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: following")
	}

	ctx := context.Background()

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feedFollows {
		fmt.Println(feed.Name)
	}

	return nil
}

func HandlerAgg(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ctx := context.Background()
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := services.ScrapeFeeds(ctx, s); err != nil {
			fmt.Println("error scraping feeds:", err)
		}
	}
}

func HandlerBrowse(s *cli.State, cmd cli.Command, user database.User) error {
	limit := int32(2)

	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: browse <limit>")
	}

	if len(cmd.Args) == 1 {
		parsed, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(parsed)
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%s\n%s\n\n", post.Title, post.Url)
	}

	return nil
}
