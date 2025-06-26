package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/henryEto/go-gator/internal/database"
)

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("no positional arguments expected, got (%d)", len(cmd.Args))
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch feeds: %w", err)
	}

	for _, feed := range feeds {
		printFeed(s, feed)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <url>", cmd.Name)
	}

	ctx := context.Background()

	feedParams := database.CreateRSSFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateRSSFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(s, feed)
	fmt.Println("=====================================")
	fmt.Println()

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return err
	}

	return nil
}

func printFeed(s *state, feed database.Feed) {
	user, err := s.db.GetUserByID(context.Background(), feed.UserID)
	if err != nil {
		log.Fatalf("error fetching the user: %v", err)
	}

	fmt.Println("===================================================")
	fmt.Printf("  Name:            %s\n", feed.Name)
	fmt.Println("---------------------------------------------------")
	fmt.Printf("  * ID:            %s\n", feed.ID)
	fmt.Printf("  * Created:       %v\n", feed.CreatedAt)
	fmt.Printf("  * Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("  * URL:           %s\n", feed.Url)
	fmt.Printf("  * User:          %s\n", user.Name)
	fmt.Println("")
}

func printFollow(follow database.FeedFollow) {
}
