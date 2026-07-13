package main

import (
	"github.com/CamdenAJohnson/Gator/internal/config"
	"github.com/CamdenAJohnson/Gator/internal/database"
)

// main instance sturct
type state struct {
	cfg *config.Config
	db *database.Queries
}

// Command sturcts
type command struct {
	name string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

// Http RSSFeed sturcts
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}