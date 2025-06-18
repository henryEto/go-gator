package main

import (
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	err := s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not set the current user: %w", err)
	}

	log.Println("User switched successfully")
	return nil
}
