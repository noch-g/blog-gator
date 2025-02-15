package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/noch-g/blog-gator/internal/database"
)

func handlerBrowse(s *state, c command, user database.User) error {
	postLimit := 2
	if len(c.Args) == 1 {
		limitArg, err := strconv.Atoi(c.Args[0])
		if err == nil {
			postLimit = limitArg
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(postLimit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedID)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
