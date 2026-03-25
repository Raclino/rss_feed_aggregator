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
	dbQueries := database.New(db)
	state := &cli.State{
		Config: conf,
		Db:     dbQueries,
	}

	commands := &cli.Commands{
		Cmd: make(map[string]func(*cli.State, cli.Command) error),
	}
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatal("error, not enough arguments provided")
	}
	if len(userArgs) == 2 {
		log.Fatal("error, username is required")
	}

	usrCommand := cli.Command{Name: userArgs[1], Args: userArgs[1:]}
	err = commands.Run(state, usrCommand)
	if err != nil {
		log.Fatal(err)
	}

}
