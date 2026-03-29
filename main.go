package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/cli/handlers"
	"github.com/Raclino/rss_feed_aggregator/internal/config"
	"github.com/Raclino/rss_feed_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	state := &cli.State{
		Config: conf,
		Db:     dbQueries,
	}

	commands := &cli.Commands{
		Handlers: make(map[string]func(*cli.State, cli.Command) error),
	}

	commands.Register("login", handlers.HandlerLogin)
	commands.Register("register", handlers.HandlerRegister)
	commands.Register("reset", handlers.HandlerReset)
	commands.Register("users", handlers.HandlerListUsers)
	commands.Register("agg", handlers.HandlerAgg)
	commands.Register("addfeed", cli.MiddlewareLoggedIn(handlers.HandlerAddFeed))
	commands.Register("feeds", handlers.HandlerListFeeds)
	commands.Register("follow", cli.MiddlewareLoggedIn(handlers.HandlerFollow))
	commands.Register("unfollow", cli.MiddlewareLoggedIn(handlers.HandlerUnFollow))
	commands.Register("following", cli.MiddlewareLoggedIn(handlers.HandlerFollowing))
	commands.Register("browse", cli.MiddlewareLoggedIn(handlers.HandlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("error: no command provided")
	}

	cmd := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := commands.Run(state, cmd); err != nil {
		log.Fatal(err)
	}
}
