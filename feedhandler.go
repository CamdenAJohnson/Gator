package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CamdenAJohnson/Gator/internal/database"
	"github.com/google/uuid"
)

func handleAddFeed(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 2 {
		return formatArgErr("addfeed <name> <url>", 2, n)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("User %v does not exist within the database: %v\n", s.cfg.Current_user_name, err)
	}
	
	if _, err := fetchFeed(context.Background(), cmd.arguments[1]); err != nil {
		return fmt.Errorf("Requset function failed: %v\n", err)
	}

	var feedParams = database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.arguments[0],
		Url: cmd.arguments[1],
		UserID: user.ID,
	}

	Feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Failed to insert RSSFeed into database: %v\n", err)
	} else {
		fmt.Println(Feed)
	}

	var followParams = database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: Feed.ID,
	}

	if _, err := s.db.CreateFeedFollow(context.Background(), followParams); err != nil {
		return fmt.Errorf("Failed to create feed follow for user: %v and feed: %v\n", user.Name, Feed.Name)
	}

	return nil
}

func handleFeeds(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 0 {
		return formatArgErr("feeds", 0, n)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrive feeds from database: %v\n", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			log.Printf("Deleted user linked to feed: %v\n", feed.ID)
			continue
		} // should never happen

		fmt.Printf("- %v\n- %v\n- %v\n\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

// Creates a feed follow record for the current user
func handleFollow(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 1 {
		return formatArgErr("follow <url>", 1, n)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Failed to retrive feed: %v\n", err)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("Failed to retive user: %v\n", err)
	}

	followParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %v\n", err)
	}

	fmt.Printf("user: %v:\n following: %v\n", follow.UserName, follow.FeedName)

	return nil
}

// Prints a list of all the feeds the current user is following
func handleFollowing(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 0 {
		return formatArgErr("following", 0, n)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("Failed to retrive user: %v\n", err)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to retive followed feeds: %v\n", err)
	}

	fmt.Printf("List of followed feeds for %v\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf(" - %v\n", feed.FeedName)
	}

	return nil
}