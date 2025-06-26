package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/henryEto/go-gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	params := database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Println("============================================================")
		fmt.Printf("%s - %v", p.Title, p.PublishedAt.Time)
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println(p.Url)
		fmt.Println(p.Description)
		fmt.Println("")
	}

	return nil
}

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
		pubTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", f.PubDate)
		if err != nil {
			pubTime = time.Now().UTC()
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       f.Title,
			Url:         f.Link,
			Description: sql.NullString{String: f.Description},
			PublishedAt: sql.NullTime{Time: pubTime},
			FeedID:      feedToFetch.ID,
		}
		err = s.db.CreatePost(context.Background(), postParams)
		if err != nil && !strings.Contains(err.Error(), "duplicate ") {
			log.Fatal(err)
		}
		// fmt.Printf("fetched '%s' from '%s'\n", f.Title, rssFeed.Channel.Title)
	}

	return nil
}
