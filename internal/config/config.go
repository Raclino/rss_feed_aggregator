package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error reading user home dir: %w", err)
	}

	return filepath.Join(homeDir, ".gatorconfig.json"), nil
}

func Read() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error Unmarshalling config: %w", err)
	}

	return &config, nil
}

func (c *Config) SetUser(currentUserName string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	c.CurrentUserName = currentUserName

	marshaledConf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, marshaledConf, 0644); err != nil {
		return fmt.Errorf("couldn't write config file: %w", err)
	}

	return nil
}
