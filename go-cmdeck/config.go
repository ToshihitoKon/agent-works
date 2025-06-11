package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ExecutionResult struct {
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
	ExitCode    int       `json:"exit_code"`
	Output      string    `json:"output,omitempty"`
}

type Context struct {
	Name         string            `json:"name"`
	Label        string            `json:"label"`
	Description  string            `json:"description,omitempty"`
	Commands     map[string]string `json:"commands"`
	Variables    map[string]string `json:"variables,omitempty"`
	LastResult   *ExecutionResult  `json:"last_result,omitempty"`
}

type ColorTheme struct {
	Title        string `json:"title"`
	Selected     string `json:"selected"`
	Border       string `json:"border"`
	OutputTitle  string `json:"output_title"`
}

type Config struct {
	Contexts map[string]Context `json:"contexts"`
	Theme    ColorTheme         `json:"theme"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "go-cmdeck", "config.json"), nil
}

func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Contexts: make(map[string]Context),
			Theme: ColorTheme{
				Title:       "205",
				Selected:    "199",
				Border:      "168",
				OutputTitle: "212",
			},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Theme.Title == "" {
		config.Theme = ColorTheme{
			Title:       "205",
			Selected:    "199",
			Border:      "168", 
			OutputTitle: "212",
		}
	}

	return &config, nil
}

func (c *Config) save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}