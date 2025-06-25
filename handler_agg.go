package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zonne13/go-gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs(1s, 1m, 1h)>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", cmd.Args[0])
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	markedParams := database.MarkFeedFetchedParams{
		ID:        feedToFetch.ID,
		UpdatedAt: time.Now().UTC(),
	}
	err = s.db.MarkFeedFetched(context.Background(), markedParams)
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	for _, f := range rssFeed.Channel.Items {
		fmt.Printf("fetched '%s' from '%s'\n", f.Title, rssFeed.Channel.Title)
	}

	return nil
}
