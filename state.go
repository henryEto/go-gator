package main

import (
	"github.com/henryEto/go-gator/internal/config"
	"github.com/henryEto/go-gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
