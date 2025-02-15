package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/google/uuid"
	"github.com/noch-g/blog-gator/internal/database"
)

func handlerFeedAggregator(s *state, c command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	fmt.Println(feed.Channel.Item)
	return nil
}

func handlerAddFeed(s *state, c command) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", c.Name)
	}
	name := c.Args[0]
	url := c.Args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully!")
	fmt.Printf("User %s now follows feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerListFeeds(s *state, c command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		fmt.Printf("%s: %s - %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}
