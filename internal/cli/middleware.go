package cli

import (
	"context"
	"fmt"

	"github.com/Raclino/rss_feed_aggregator/internal/database"
)

func MiddlewareLoggedIn(
	handler func(s *State, cmd Command, user database.User) error,
) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		if s.Config.CurrentUserName == "" {
			return fmt.Errorf("no user is currently logged in")
		}
		ctx := context.Background()

		user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
