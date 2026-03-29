package handlers

import (
	"context"
	"fmt"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
)

func HandlerReset(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	if err := s.Db.DeleteAllUsers(ctx); err != nil {
		return err
	}

	fmt.Println("Database Reset was successful")
	return nil
}
