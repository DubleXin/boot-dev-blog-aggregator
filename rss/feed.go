package rss

import "html"

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
