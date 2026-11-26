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

func handlerAddFeed(s *state, cmd command, user database.User) error {

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

	feedParams := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := addFeed(s, feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	if err = createFeedFollow(s, feedFollowParams); err != nil {
		return fmt.Errorf("Failed to create feed-follow record, with given arguments: user-id: %s \t feed-id: %d\n\tError: %s", feed.UserID.String(), feed.ID, err.Error())
	}

	fmt.Printf("Successfully added feed \tname: \"%s\"\tlink: \"%s\"\n", feed.Name, feed.Url)
	return nil
}

func addFeed(s *state, params database.CreateFeedParams) (database.Feed, error) {
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return database.Feed{}, fmt.Errorf("Failed to create feed. Error:\n\t\t%s", err.Error())
	}
	return feed, nil
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

func createFeedFollow(s *state, params database.CreateFeedFollowParams) error {
	if _, err := s.db.CreateFeedFollow(context.Background(), params); err != nil {
		return fmt.Errorf("Failed to create feed-follow record")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 1 {
		return fmt.Errorf("The follow handler expects single argument, feed url")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to get feed by url: \"%s\"\n\t\tError: %s", url, err.Error())
	}

	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	if err = createFeedFollow(s, params); err != nil {
		return fmt.Errorf("Failed to create feed-follow record, with given arguments: user-id: %s \t feed-id: %d\n\tError: %s", feed.UserID.String(), feed.ID, err.Error())
	}

	fmt.Printf("Successfully added feed-follow record: \"%s\" : \"%s\"\n", feed.Name, feed.Url)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	records, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get records for user with id: %s", user.ID.String())
	}
	fmt.Println("Following list:")
	for _, record := range records {
		fmt.Printf("* name: %s\turl: %s\n", record.Name, record.Url)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 1 {
		return fmt.Errorf("The follow handler expects single argument, feed url")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to get feed by url: \"%s\"\n\t\tError: %s", url, err.Error())
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.db.DeleteFeedFollow(context.Background(), params); err != nil {
		return fmt.Errorf("Failed ro delete feed-follow record.\n\tError: %s", err.Error())
	}

	fmt.Printf("Successfully deleted feed-follow record")
	return nil
}
