package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetDB(context.Background())
	if err != nil {
		log.Fatalf("couldn't reset database: %v", err)
		return err
	}
	fmt.Println("Database was reset successfully")
	return nil
}
