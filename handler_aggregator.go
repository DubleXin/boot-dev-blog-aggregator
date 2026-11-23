package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DubleXin/aggregator/internal/database"
	"github.com/DubleXin/aggregator/rss"
)

func getFeed() (*rss.RSSFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
}

func handlerAggregator(s *state, cmd command) error {
	feed, err := getFeed()

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return err
	}

	if len(cmd.args) < 2 {
		return fmt.Errorf("The add feed handler expects two arguments, feed name and feed url.\n\t Example: addfeed \"Hacker News RSS\" \"https://hnrss.org/newest\"")
	}
	name, url := cmd.args[0], cmd.args[1]

	if len(name) == 0 {
		return fmt.Errorf("Cannot add feed user with empty name.")
	}
	if len(url) == 0 {
		return fmt.Errorf("Cannot add feed user with empty url.")
	}

	params := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create feed. Error:\n\t\t%s", err.Error())
	}

	fmt.Printf("Successfully added feed \tname: \"%s\"\tlink: \"%s\"\n", feed.Name, feed.Url)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedUsersList(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds. Error:\n\t\t%w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("* feed name/url: \"%s\"/\"%s\"\tusername: \"%s\"\n", feed.Name, feed.Url, feed.Username)
	}
	return nil
}
