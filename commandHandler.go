package main

import (
	"context"
	"fmt"
	"log"

	"github.com/CamdenAJohnson/Gator/internal/database"
)

// Registers the give coomand function
func (c *commands) register(name string, f func(*state, command) error) {
	_, exist := c.handlers[name]
	if exist == true {
		log.Printf("A command with the name %v alread exist", name)
	}

	c.handlers[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	f := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil || user.Name != s.cfg.Current_user_name {
			return fmt.Errorf("User check failed: %v\n", err)
		}

		return handler(s, cmd, user)
	}

	return f
}

func formatArgErr(usage string, expected, got int) error {
	fmt.Printf("Usage: %v\n", usage)
	return fmt.Errorf("Incorrect number of arguments.\n - Expected: %v\n - Got: %v\n", expected, got)
}

func handleAgg(s *state, cmd command) error {
	// WIP
	Feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Request function failed: %v\n", err)
	}

	fmt.Println(Feed)
	return nil
}

