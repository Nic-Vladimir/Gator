package main

import (
	"context"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Unfollow requires feed URL")
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Couldn't find feed: %w", err)
	}
	unfollowFeedParams := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	_, err = s.db.UnfollowFeed(context.Background(), unfollowFeedParams)
	if err != nil {
		return fmt.Errorf("Couldn't unfollow feed: %w", err)
	}
	fmt.Printf("User %s unfollowed %s\n", user.Name, feed.Name)
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Follow requires feed name")
	}
	feedUrl := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("Couldn't add feed follow: %w", err)
	}
	fmt.Printf("User %s followed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerListFollowing(s *state, cmd command, user database.User) error {
	followingList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Couldn't get following feeds: %w", err)
	}
	for _, feedFollow := range followingList {
		feed, err := s.db.GetFeedByID(context.Background(), feedFollow.FeedID)
		if err != nil {
			return fmt.Errorf("Couldn't get feed: %w", err)
		}
		fmt.Println(" *", feed.Name)
	}
	return nil
}
