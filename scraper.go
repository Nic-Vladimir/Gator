package main

import (
	"context"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

func scrapeFeeds(s *state) error {
	feedToScrape, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't get feed to scrape: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), feedToScrape.ID)
	if err != nil {
		return fmt.Errorf("Couldn't mark feed as fetched: %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), feedToScrape.Url)
	if err != nil {
		return fmt.Errorf("Couldn't fetch feed: %w", err)
	}
	fmt.Printf("Fetched %d items from %s\n", len(rssFeed.Channel.Item), feedToScrape.Url)
	for _, item := range rssFeed.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		createPostsParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feedToScrape.ID,
		}
		inserted, err := s.db.CreatePost(context.Background(), createPostsParams)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				fmt.Println("Skipping duplicate post")
				continue
			}
			fmt.Printf("Couldn't create post %s: %s", inserted.Title, err)
			continue
		}
		fmt.Printf("Created post: %s\n", inserted.Title)
	}
	return nil
}
