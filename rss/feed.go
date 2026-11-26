package rss

import (
	"fmt"
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (r *RSSFeed) unescapeSelf() {
	if r.Channel.Title != "" {
		r.Channel.Title = html.UnescapeString(r.Channel.Title)
	}
	if r.Channel.Description != "" {
		r.Channel.Description = html.UnescapeString(r.Channel.Description)
	}
	for _, item := range r.Channel.Item {
		if item.Title != "" {
			item.Title = html.UnescapeString(item.Title)
		}
		if item.Description != "" {
			item.Description = html.UnescapeString(item.Description)
		}
	}
}

func (feed *RSSFeed) Print() {
	if feed == nil {
		fmt.Println("RSS feed is empty.")
		return
	}

	fmt.Printf("Title: %s\n", feed.Channel.Title)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Items: %d\n", len(feed.Channel.Item))

	for i, item := range feed.Channel.Item {
		fmt.Printf("\n--- Item %d ---\n", i+1)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Printf("Published: %s\n", item.PubDate)
	}
}
