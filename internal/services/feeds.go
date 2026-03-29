package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/Raclino/rss_feed_aggregator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var ErrDupKeyUrl = errors.New("duplicate key value violates unique constraint 'posts_url_key'")

func ScrapeFeeds(ctx context.Context, s *cli.State) error {
	nxtFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	fmt.Println(nxtFeed)

	now := time.Now()
	markFeedFetchedParams := database.MarkFeedFetchedParams{
		ID:        nxtFeed.ID,
		UpdatedAt: now,
		LastFetchedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}
	f, err := s.Db.MarkFeedFetched(ctx, markFeedFetchedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching feed: %s\n", nxtFeed.Name)
	rssFeed, err := rss.FetchFeed(ctx, f.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		parsedTime, err := parsePublishedAt(item.PubDate)
		if err != nil {
			return fmt.Errorf("couldn't parse publishedAt time format %w", err)
		}
		now := time.Now()

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsedTime,
			FeedID:      nxtFeed.ID,
		}
		post, err := s.Db.CreatePost(ctx, postParams)
		if err != nil {
			// duplicate key, ignored
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			fmt.Printf("error creating post: %v\n", err)
		}
		if err != nil {
			return err
		}
		fmt.Printf("created Post: %+v\n", post)
	}
	return nil
}
func parsePublishedAt(s string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{Valid: false}, nil
	}

	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return sql.NullTime{
				Time:  t,
				Valid: true,
			}, nil
		}
	}

	return sql.NullTime{}, fmt.Errorf("could not parse published_at: %q", s)
}
