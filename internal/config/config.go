package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Error reading user homeDir: %w", err)
	}

	gatorFile, err := os.OpenInRoot(homeDir, ".gatorconfig.json")
	if err != nil {
		return nil, fmt.Errorf("Error opening .gatorconfig.json file: %w", err)
	}
	data, err := io.ReadAll(gatorFile)
	if err != nil {
		return nil, fmt.Errorf("Error while reading gatorFile: %w", err)
	}

	config := Config{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Error Unmarshalling data into Config struct: %w", err)
	}

	return &config, nil
}

func (c *Config) SetUser(currentUserName string) error {

	return nil
}
