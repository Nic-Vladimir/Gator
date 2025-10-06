package main

import (
	"context"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
	"strconv"
)

func handlerBrowseFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("browse [post_limit:int]")
	}
	limit := 2
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("browse [post_limit:int] ")
		}
		limit = parsedLimit
	}
	getPostsForUserParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostsForUserParams)
	if err != nil {
		return fmt.Errorf("Couldn't get posts for %s: %w", user.Name, err)
	}
	for _, p := range posts {
		fmt.Println(p.Title)
		fmt.Println(p.Url)
		fmt.Println(p.PublishedAt)
		fmt.Println()
	}
	return nil
}
