package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the req: %w", err)
	}

	req.Header.Set("User-agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the req: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read from the body: %w", err)
	}
	rssFeed := RSSFeed{}

	if err = xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, fmt.Errorf("couldn't Unmarshal the response: %w", err)
	}
	html.UnescapeString(rssFeed.Channel.Title)
	html.UnescapeString(rssFeed.Channel.Description)
	fmt.Println(rssFeed)

	return &rssFeed, nil
}
