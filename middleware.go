package main

import (
	"context"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// basic sanity checks
		if s == nil {
			return fmt.Errorf("internal error: nil state")
		}
		if s.config == nil {
			return fmt.Errorf("internal error: missing config")
		}
		if s.db == nil {
			return fmt.Errorf("internal error: missing db")
		}

		username := s.config.CurrentUserName
		if username == "" {
			return fmt.Errorf("not logged in")
		}

		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("login required: %w", err)
		}

		return handler(s, cmd, user)
	}
}
