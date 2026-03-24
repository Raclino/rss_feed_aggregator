package main

import (
	"fmt"
	"log"

	"github.com/Raclino/rss_feed_aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf)
}
