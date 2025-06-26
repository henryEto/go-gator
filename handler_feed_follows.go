package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/henryEto/go-gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed does not exist: %w", err)
	}

	unfollowParam := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.Unfollow(context.Background(), unfollowParam)
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("%s no longer follows %s", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("no agruments expected: <%d> given", len(cmd.Args))
	}

	folows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
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

func handlerFollowFeed(s *state, cmd command, user database.User) error {
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
