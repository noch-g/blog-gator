package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/noch-g/blog-gator/internal/database"
)

func handlerFollow(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_name>", c.Name)
	}
	feedURL := c.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't get feed %s: %w", feedURL, err)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user %s: %w", s.cfg.CurrentUserName, err)
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

	fmt.Printf("User %s now follows feed %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, c command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user %s: %w", s.cfg.CurrentUserName, err)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows for user %s: %w", user.Name, err)
	}

	if len(feedFollows) == 0 {
		fmt.Printf("User %s doesn't follow any feed\n", user.Name)
		return nil
	}

	fmt.Printf("User %s follows %d feeds\n", user.Name, len(feedFollows))

	for _, feedFollow := range feedFollows {
		feed, err := s.db.GetFeedByID(context.Background(), feedFollow.FeedID)
		if err != nil {
			return fmt.Errorf("couldn't get feed %s: %w", feedFollow.FeedID, err)
		}
		fmt.Printf("%s: %s\n", feed.Name, feed.Url)
	}

	return nil
}
