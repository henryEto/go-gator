package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/henryEto/go-gator/internal/config"
	"github.com/henryEto/go-gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
		return
	}
	dbQueries := database.New(db)

	programState := &state{
		cfg: cfg,
		db:  dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", midlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", midlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", midlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", midlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", midlewareLoggedIn(handlerBrowse))
	cmds.register("help", handlerHelp)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

func midlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Username)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func handlerHelp(s *state, cmd command) error {
	fmt.Println("Welcome to go-gator! Here are the available commands:")
	fmt.Println("----------------------------------------------------")

	// Print help for each command. This can be extended with more detailed descriptions
	// if we add them to the command registration.
	// For now, we'll list them with their basic usage patterns.
	fmt.Println("  login <username>    - Log in as an existing user.")
	fmt.Println("  register <username> - Register a new user.")
	fmt.Println("  reset               - Resets the database (use with caution!).")
	fmt.Println("  users               - List all registered users.")
	fmt.Println("  agg <time>          - Start aggregation of feeds (e.g., 'agg 1m', 'agg 5s').")
	fmt.Println("  addfeed <name> <url>- Add a new RSS feed and automatically follow it.")
	fmt.Println("  feeds               - List all available RSS feeds.")
	fmt.Println("  follow <url>        - Follow an existing RSS feed by its URL.")
	fmt.Println("  following           - List all feeds currently followed by the logged-in user.")
	fmt.Println("  unfollow <url>      - Unfollow a previously followed RSS feed.")
	fmt.Println("  browse [limit]      - Browse posts from your followed feeds (optional limit, default 2).")
	fmt.Println("  help                - Display this help message.")
	fmt.Println("----------------------------------------------------")
	fmt.Println("To run a command: go run . <command_name> [arguments]")
	fmt.Println("Example: go run . register myuser")

	return nil
}
