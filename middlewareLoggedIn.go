package main

import (
	"context"
	"fmt"

	"github.com/DubleXin/aggregator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return fmt.Errorf("Denied. User should be logged in.\n%s", err.Error())
		}

		return handler(s, cmd, user)
	}
}
