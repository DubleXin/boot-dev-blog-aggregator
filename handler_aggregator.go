package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DubleXin/aggregator/internal/database"
	"github.com/DubleXin/aggregator/rss"
)

func getFeed(url string) (*rss.RSSFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return rss.FetchFeed(ctx, url)
}

func scrapeFeeds(s *state, user database.User) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx, user.ID)
	if err != nil {
		return fmt.Errorf(
			"failed to get next feed in queue: %s",
			err.Error(),
		)
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to mark feed: %s", err.Error())
	}

	content, err := getFeed(feed.Url)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch feed %q: %s",
			feed.Url,
			err.Error(),
		)
	}

	for _, item := range content.Channel.Item {
		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			println(
				"Failed to parse publication date for post %q: %s",
				item.Link,
				err.Error(),
			)
		}

		now := time.Now()

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			CreatedAt: now,
			UpdatedAt: now,
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			println(
				"Failed to save post %q from feed %q: %s",
				item.Link,
				feed.Name,
				err.Error(),
			)
		}
	}

	return nil
}

func parsePublishedAt(value string) (sql.NullTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return sql.NullTime{}, nil
	}

	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339Nano,
		time.RFC850,
		time.ANSIC,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 -0700",
		"2006-01-02 15:04:05 -0700",
	}

	var lastErr error

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return sql.NullTime{
				Time:  parsed.UTC(),
				Valid: true,
			}, nil
		}

		lastErr = err
	}

	return sql.NullTime{}, fmt.Errorf(
		"unsupported publication date %q: %s",
		value,
		lastErr.Error(),
	)
}

func handlerAggregator(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 1 {
		return fmt.Errorf("The aggregator handler expects single time betwwen requests")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])

	if err != nil {
		return fmt.Errorf("Failed to parse duration, for input: %s", cmd.args[0])
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf(
			"browse expects at most one argument.\n\tExample: browse 5",
		)
	}

	limit := int32(2)

	if len(cmd.args) == 1 {
		parsed, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid limit: %s", err.Error())
		}
		if parsed <= 0 {
			return fmt.Errorf("limit must be greater than zero")
		}

		limit = int32(parsed)
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  limit,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get posts: %s", err.Error())
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	for i, post := range posts {
		fmt.Printf("\n--- Post %d ---\n", i+1)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)

		if post.PublishedAt.Valid {
			fmt.Printf(
				"Published: %s\n",
				post.PublishedAt.Time.Format("2006-01-02 15:04:05 MST"),
			)
		}

		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
	}

	return nil
}
