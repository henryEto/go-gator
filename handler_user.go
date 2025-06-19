package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zonne13/go-gator/internal/database"
)

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.Username {
			fmt.Printf("  * %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("  * %v\n", user.Name)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	userDB, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not get user from database: %s", err)
	}

	err = s.cfg.SetUser(userDB.Name)
	if err != nil {
		return fmt.Errorf("could not set the current user: %w", err)
	}

	fmt.Println("User switched successfully")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Errorf("could not create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}
	fmt.Println("User switched successfully")
	return nil
}
