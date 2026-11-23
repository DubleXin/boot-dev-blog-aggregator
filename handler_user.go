package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/DubleXin/aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("The login handler expects a single argument, the username")
	}
	username := cmd.args[0]

	if len(username) == 0 {
		return fmt.Errorf("Cannot login user with empty username.")
	}

	user, err := s.db.GetUser(context.Background(), username)

	if err != nil {
		return fmt.Errorf("User \"%s\" does not exist.", username)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Failed to set user inside config. Error:\n\t\t%s", err.Error())
	}

	fmt.Printf("Successfully logged in as \"%s\"\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("The register handler expects a single argument, the username")
	}
	name := cmd.args[0]
	if len(name) == 0 {
		return fmt.Errorf("Cannot register user with empty username.")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), params)

	if err != nil {
		return fmt.Errorf("Failed to create user. Error:\n\t\t%s", err.Error())
	}

	err = s.cfg.SetUser(user.Name)

	if err != nil {
		return fmt.Errorf("Failed to set user inside config. Error:\n\t\t%s", err.Error())
	}

	fmt.Printf("Successfully registered user \"%s\"\n", user.Name)
	return nil
}

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users table")
	}
	fmt.Println("Successfully reset user table")
	return nil
}

func getUsers(s *state) ([]string, error) {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := getUsers(s)
	if err != nil {
		return err
	}

	for _, username := range users {
		fmt.Printf("* %s", username)
		if username == s.cfg.CurrentUser {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}
