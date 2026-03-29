package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	"github.com/Raclino/rss_feed_aggregator/internal/rss"
)

func ScrapeFeeds(ctx context.Context, s *cli.State) error {
	nxtFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	fmt.Println(nxtFeed)

	markFeedFetchedParams := database.MarkFeedFetchedParams{
		ID:        nxtFeed.ID,
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	f, err := s.Db.MarkFeedFetched(ctx, markFeedFetchedParams)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(ctx, f.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
