package main

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"time"

	"github.com/google/uuid"
	"github.com/noch-g/blog-gator/internal/database"
)

func handlerFeedAggregator(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>, example: %s 10s", c.Name, c.Name)
	}
	timeBetweenReqs := c.Args[0]

	timeBetweenReqsDuration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("couldn't parse time_between_reqs: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqsDuration)

	ticker := time.NewTicker(timeBetweenReqsDuration)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %s\n", err)
		}
	}
}

func handlerAddFeed(s *state, c command, user database.User) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", c.Name)
	}
	name := c.Args[0]
	url := c.Args[1]

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

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	s.db.MarkFeedFetched(context.Background(), feed.ID)

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	items := fetchedFeed.Channel.Item
	postCount := 0
	for _, item := range items {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     html.UnescapeString(item.Title),
			Url:       item.Link,
			Description: sql.NullString{
				String: html.UnescapeString(item.Description),
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				continue
			}
			fmt.Printf("couldn't create post: %v\n", err)
			continue
		}
		postCount++
	}
	fmt.Printf("Found %d new posts from %s\n", postCount, feed.Url)

	return nil
}
