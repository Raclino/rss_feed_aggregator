package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
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
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerUsers)
	commands.Register("agg", cli.HandlerAgg)
	commands.Register("addfeed", cli.HandlerAddFeed)
	commands.Register("feeds", cli.HandlerFeeds)

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
