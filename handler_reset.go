package main

import (
	"context"
	"fmt"
)

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't delete users: %w", err)
	}
	fmt.Println("Users table reset")
	return nil
}
