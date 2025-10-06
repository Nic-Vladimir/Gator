package main

import (
	"context"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't list users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf(" * %s (current)\n", user.Name)
			continue
		}
		fmt.Println(" *", user.Name)
	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" ID:	%v\n", user.ID)
	fmt.Printf(" Name:	%v\n", user.Name)
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Register requires username")
	}
	username := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	s.config.SetUser(username)
	fmt.Printf("Registered user: %s\n", username)
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Login requires username")
	}
	username := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Login failed: user '%s' does not exist", username)
	}

	err = s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("Logged in as: %s\n", cmd.Args[0])
	return nil
}
