package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CamdenAJohnson/Gator/internal/database"
	"github.com/google/uuid"
)

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