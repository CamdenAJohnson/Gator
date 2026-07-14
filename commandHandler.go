package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

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

func handleAgg(s *state, cmd command, user database.User) error {
	if n := len(cmd.arguments); n != 1 {
		return formatArgErr("agg <interval>", 1, n)
	}

	num, err := strconv.Atoi(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Argument is not a number: %v\n", err)
	}

	if num <= 0 { num = 1 }
	interval := time.Duration(num) * time.Minute
	ticker := time.NewTicker(interval)

	fmt.Printf("time between reqs: %v\n", interval)

	for range ticker.C {
		if err := aggHelper(s, user); err != nil {
			log.Printf("Error during aggregation: %v\n", err)
		}
	}
	
	return nil
}

func aggHelper(s *state, user database.User) error {
	oldestFeed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to retrive oldest feed: %v\n", err)
	}

	feed, err := fetchFeed(context.Background(), oldestFeed.Url)
	if err != nil {
		return fmt.Errorf("Fetch request function failed: %v\n", err)
	}

	validTime := sql.NullTime{
		Time: time.Now(),
		Valid: true,
	}

	markFeedParams := database.MarkFeedFetchedParams{
		LastFetchedAt: validTime,
		UpdatedAt: time.Now(),
		ID: oldestFeed.ID,
	}

	if err := s.db.MarkFeedFetched(context.Background(), markFeedParams); err != nil {
		return fmt.Errorf("Failed to update fetched feed: %v\n", err)
	}

	fmt.Printf("Title: %v\nDescription: %v\n", feed.Channel.Title, feed.Channel.Description)
	for n, elem := range feed.Channel.Item {
		fmt.Printf("Item: %v\nTitle: %v\nDescription: %v\n", n, elem.Title, elem.Description)
	}

	fmt.Printf("\n\n-----------------------\n\n")

	return nil
}