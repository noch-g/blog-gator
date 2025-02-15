package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/noch-g/blog-gator/internal/database"
)

func handlerFollow(s *state, c command, user database.User) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", c.Name)
	}
	feedURL := c.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't get feed %s: %w", feedURL, err)
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

func handlerFollowing(s *state, c command, user database.User) error {
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

func handlerUnfollow(s *state, c command, user database.User) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", c.Name)
	}
	feedURL := c.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't get feed %s: %w", feedURL, err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}

	fmt.Printf("User %s unfollowed feed %s\n", user.Name, feed.Name)

	return nil
}
