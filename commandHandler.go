package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

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

/* WIP
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	
}
*/

/* WIP
func handlerAddFeed(s *state, cmd command, user database.User) error {

}
*/

func formatArgErr(usage string, expected, got int) error {
	fmt.Printf("Usage: %v\n", usage)
	return fmt.Errorf("Incorrect number of arguments.\n - Expected: %v\n - Got: %v\n", expected, got)
}

// Command. Handles uesr login
func handleLogin(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 1 {
		return formatArgErr("login <name>", 1, n)
	}

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0]);
	if  err != nil {
		return fmt.Errorf("User does not exist: %v\n", err)
	}

	s.cfg.SetUser(user.Name)

	if s.cfg.Current_user_name != user.Name {
		return fmt.Errorf("Failed to set user: %v\n", user.Name)
	}

	fmt.Printf("User: %v has been set\n", user.Name)
	
	return nil
}

// Command Handles registering a new users to the database
func handleRegister(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 1 {
		return formatArgErr("register <name>", 1, n)
	}

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.arguments[0],
	}
	
	if _, err := s.db.GetUser(context.Background(), userParams.Name); err == nil {
		return fmt.Errorf("User already exist: %v\n", err)
	}

	user, err := s.db.CreateUser(context.Background(), userParams)

	if err != nil {
		return fmt.Errorf("Failed to create new user: %v\n", cmd.arguments[0])
	}
	
	if err := handleLogin(s, cmd); err != nil {
		return err
	}

	fmt.Printf("User: %v has been created!\n", user.Name)

	return nil
}

// Command. Handles deleting every user from the users table
func handleReset(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 0 {
		return formatArgErr("reset", 0, n)
	}

	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users table: %v\n", err)
	}

	fmt.Println("users table has been reset!")

	return nil
}

func handleUsers(s *state, cmd command) error {
	if n := len(cmd.arguments); n != 0 {
		return formatArgErr("users", 0, n)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retive users from database: %v\n", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.Current_user_name {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
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