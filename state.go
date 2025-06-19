package main

import (
	"github.com/zonne13/go-gator/internal/config"
	"github.com/zonne13/go-gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
