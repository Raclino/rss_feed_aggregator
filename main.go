package main

import (
	"log"
	"os"

	"github.com/Raclino/rss_feed_aggregator/internal/cli"
	"github.com/Raclino/rss_feed_aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	state := &cli.State{
		Config: conf,
	}

	commands := &cli.Commands{
		Cmd: make(map[string]func(*cli.State, cli.Command) error),
	}
	commands.Register("login", cli.HandlerLogin)

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatal("error, not enough arguments provided")
	}
	if len(userArgs) == 2 {
		log.Fatal("error, username is required")
	}

	usrCommand := cli.Command{Name: userArgs[1], Args: userArgs[2:]}
	err = commands.Run(state, usrCommand)
	if err != nil {
		log.Fatal(err)
	}

}
