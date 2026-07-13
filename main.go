package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/CamdenAJohnson/Gator/internal/config"
	"github.com/CamdenAJohnson/Gator/internal/database"
	_ "github.com/lib/pq"
)

var handles = commands{
	handlers: map[string]func(*state, command) error{
		"login": handleLogin,
		"register": handleRegister,
		"reset": handleReset,
		"users": handleUsers,
		"agg": handleAgg,
		"addfeed": handleAddFeed,
		"feeds": handleFeeds,
		"follow": handleFollow,
		"following": handleFollowing,
	},
}

func main() {
	var instance state
	cfg := config.Read()
	instance.cfg = &cfg

	db, err := sql.Open("postgres", instance.cfg.Db_url)
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}
	
	dbQueries := database.New(db)
	instance.db = dbQueries

	var cmds commands
	cmds.handlers = handles.handlers

	if n := len(os.Args); n < 2 {
		log.Fatalf("To few arguments")
	}

	cmd := command{
		name: os.Args[1],
		arguments: os.Args[2:],
	}
	
	if err := cmds.run(&instance, cmd); err != nil {
		log.Fatal(err)
	}
}

// Runs the given command
func (c *commands) run(s *state, cmd command) error {
	elem, exist := c.handlers[cmd.name]
	if exist == false {
		return fmt.Errorf("Command not found: %v", cmd.name)
	}

	if err := elem(s, cmd); err != nil {
		return fmt.Errorf("Command failed to exacute: %v", err)
	}

	return nil
}