package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sdhornet/gator/internal/database"
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func handlerAgg(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", *feed)

	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var f RSSFeed
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feed url said %s", resp.Status)
	}

	if err := xml.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("decoding feed: %w", err)
	}
	unescape(&f)

	return &f, nil
}

func unescape(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("requires two arguments to add a feed <name> <url>")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	now := time.Now()

	p := database.CreateFeedParams{
		ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Name: name, Url: url, UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), p)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("could not add follow row: %w", err)
	}

	fmt.Printf("Feed Name: %s\n", feed.Name)
	fmt.Printf("Feed URL: %s\n", feed.Url)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("feeds takes no arguments")
	}
	feedData, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	for _, d := range feedData {
		fmt.Printf("Name: %s\n", d.Name)
		fmt.Printf("URL: %s\n", d.Url)
		fmt.Printf("User: %s\n", d.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow requires only a URL parameter")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("feed not found for this url: %w", err)
	}

	now := time.Now()
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("could not add follow row: %w", err)
	}
	fmt.Printf("Feed:%s User:%s\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return errors.New("following takes no arguments")
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not fetch user's follows: %w", err)
	}

	for _, f := range following {
		fmt.Printf("%s\n", f.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("unfollow requires a url argument only")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed not found for this url: %w", err)
	}

	if err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{FeedID: feed.ID, UserID: user.ID}); err != nil {
		return fmt.Errorf("count not unfollow the %s feed: %w", url, err)
	}

	return nil
}
