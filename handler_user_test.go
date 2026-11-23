package main

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"testing"

	"github.com/DubleXin/aggregator/internal/config"
	"github.com/DubleXin/aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestHandlerUserLifecycle(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() failed: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		t.Fatalf("sql.Open() failed: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("db.Close() failed: %v", err)
		}
	})

	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping() failed: %v", err)
	}

	originalUsername := cfg.CurrentUser
	t.Cleanup(func() {
		if err := cfg.SetUser(originalUsername); err != nil {
			t.Errorf("Failed to restore config user: %v", err)
		}
	})

	s := &state{
		cfg: cfg,
		db:  database.New(db),
	}

	ctx := context.Background()
	username := "test_user_" + uuid.NewString()
	cmd := command{args: []string{username}}

	t.Cleanup(func() {
		if err := handlerResetUsers(s, command{}); err != nil {
			t.Errorf("handlerResetUsers() failed: %v", err)
			return
		}

		_, err := s.db.GetUser(ctx, username)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf(
				"GetUser() after reset: got error %v, want sql.ErrNoRows",
				err,
			)
		}
	})

	if err := handlerRegister(s, cmd); err != nil {
		t.Fatalf("handlerRegister() failed: %v", err)
	}

	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		t.Fatalf("GetUser() after registration failed: %v", err)
	}

	if user.Name != username {
		t.Errorf("Registered user name = %q, want %q", user.Name, username)
	}

	if cfg.CurrentUser != username {
		t.Errorf(
			"CurrentUser after registration = %q, want %q",
			cfg.CurrentUser,
			username,
		)
	}

	if err := cfg.SetUser(""); err != nil {
		t.Fatalf("Failed to clear current user before login: %v", err)
	}

	if err := handlerLogin(s, cmd); err != nil {
		t.Fatalf("handlerLogin() failed: %v", err)
	}

	if cfg.CurrentUser != username {
		t.Errorf(
			"CurrentUser after login = %q, want %q",
			cfg.CurrentUser,
			username,
		)
	}
}

func TestGetUsers(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() failed: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		t.Fatalf("sql.Open() failed: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("db.Close() failed: %v", err)
		}
	})

	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping() failed: %v", err)
	}

	s := &state{
		cfg: cfg,
		db:  database.New(db),
	}

	originalUsername := cfg.CurrentUser
	t.Cleanup(func() {
		if err := handlerResetUsers(s, command{}); err != nil {
			t.Errorf("handlerResetUsers() failed: %v", err)
		}
		if err := cfg.SetUser(originalUsername); err != nil {
			t.Errorf("Failed to restore config user: %v", err)
		}
	})

	if err := handlerResetUsers(s, command{}); err != nil {
		t.Fatalf("Initial reset failed: %v", err)
	}

	users, err := getUsers(s)
	if err != nil {
		t.Fatalf("getUsers() on empty database failed: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("Expected empty users list, got %v", users)
	}

	want := []string{"alice", "bob"}
	for _, username := range want {
		cmd := command{args: []string{username}}
		if err := handlerRegister(s, cmd); err != nil {
			t.Fatalf("handlerRegister(%q) failed: %v", username, err)
		}
	}
	got, err := getUsers(s)
	if err != nil {
		t.Fatalf("getUsers() failed: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("Users = %v, want %v", got, want)
	}

	sort.Strings(got)
	sort.Strings(want)

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Users[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
