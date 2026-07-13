package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)


type Config struct {
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to retrive user home dir: %v", err)
	}

	filePath := filepath.Join(home, configFileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config fil: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		log.Fatalf("Failed to parse json data: %v", err)
	}

	return cfg
}

func (c *Config) SetUser(current_user_name string) {
	c.Current_user_name = current_user_name
	if err := write(*c); err != nil {
		fmt.Printf("Failed to set user name: %v", err)
	}
}

func write(cfg Config) error {
    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to retrieve user home dir: %v", err)
    }

    // Best practice: Use filepath.Join instead of hardcoding "/" (works on Windows too)
    filePath := filepath.Join(home, configFileName)

    // Optional but recommended: MarshalIndent formats the JSON with line breaks and tabs
    jsonData, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal config struct: %v", err)
    }

    // os.WriteFile automatically handles truncation and closing
    if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
        return fmt.Errorf("failed to write json data to file: %v", err)
    }

    return nil
}