package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	if feedURL == "" {
		return nil, fmt.Errorf("Feed url cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Request with Context")
	}
	client := http.Client{}

	req.Header.Set("User-Agent", "Gator RSS Client")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get response. \"GET %s\"", feedURL)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code, want \"%d\", got \"%d\". \"GET %s\"", http.StatusOK, res.StatusCode, feedURL)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body bytes. \"GET %s\"", feedURL)
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal xml response. \"GET %s\"", feedURL)
	}
	feed.unescapeSelf()

	return &feed, nil
}
