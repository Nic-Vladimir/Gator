package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't list feeds: %w", err)
	}
	for _, f := range feeds {
		addedBy, err := s.db.GetUserByID(context.Background(), f.UserID)
		if err != nil {
			return fmt.Errorf("Couldn't find user who added %s: %w", f.Name, err)
		}
		fmt.Println("Name: ", f.Name)
		fmt.Println("Url: ", f.Url)
		fmt.Println("Added by: ", addedBy.Name)
		fmt.Println()
	}
	return nil
}

func handlerAggregateRss(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("agg requires <time_between_reqs>")
	}
	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return errors.New("input a duration in the format: 1h2m3s")
	}
	fmt.Println("Collecting feeds every", interval)
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("addfeed requires <name> <url>")
	}
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	err = handlerFollowFeed(s, command{Name: "follow", Args: []string{cmd.Args[1]}}, user)
	if err != nil {
		return err
	}
	fmt.Printf("Created feed: %s\n", feed.Name)
	fmt.Printf("ID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\nUrl: %v\nUserID: %v",
		feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)

	return nil
}
