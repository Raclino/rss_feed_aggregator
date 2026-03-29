package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/Raclino/rss_feed_aggregator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func ScrapeFeeds(ctx context.Context, s *cli.State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	feed, err := s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:        nextFeed.ID,
		UpdatedAt: now,
		LastFetchedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Fetching feed: %s\n", feed.Name)

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			fmt.Printf("error parsing date for post %q: %v\n", item.Title, err)
			publishedAt = sql.NullTime{Valid: false}
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		_, err = s.Db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			fmt.Printf("error creating post %q: %v\n", item.Title, err)
		}
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
