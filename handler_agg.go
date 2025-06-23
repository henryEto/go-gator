package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/zonne13/go-gator/internal/database"
)

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("no agruments expected: <%d> given", len(cmd.Args))
	}

	folows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.Username)
	if err != nil {
		return err
	}
	if len(folows) == 0 {
		fmt.Printf("%s doesn't follow any feeds yet...", s.cfg.Username)
		return nil
	}
	printFollows(folows)
	return nil
}

func printFollows(follows []database.GetFeedFollowsForUserRow) {
	fmt.Printf("n%s follows %d feed(s):\n", follows[0].UserName, len(follows))
	for i, follow := range follows {
		fmt.Printf("  %d. %s\n", i+1, follow.FeedName)
	}
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	var feed database.Feed
	feedExists := false
	for _, f := range feeds {
		if f.Url == cmd.Args[0] {
			feed = f
			feedExists = true
			break
		}
	}
	if !feedExists {
		return errors.New("feed does not exist")
	}

	user, err := s.db.GetUser(ctx, s.cfg.Username)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	follow, err := s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return err
	}

	fmt.Printf("%s now folows %s!\n", follow.UserName, follow.FeedName)

	return nil
}

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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <url>", cmd.Name)
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

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

func handlerAgg(s *state, cmd command) error {
	laneRssURL := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), laneRssURL)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %v\n", rssFeed)
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
