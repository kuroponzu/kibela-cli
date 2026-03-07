package config

import (
	"fmt"
	"os"
)

// Config holds the configuration for Kibela API access.
type Config struct {
	Token string
	Team  string
}

// ExitCode represents different error types.
const (
	ExitOK           = 0
	ExitConfig       = 1
	ExitAuth         = 2
	ExitPermission   = 3
	ExitNotFound     = 4
	ExitInput        = 10
	ExitGraphQL      = 20
)

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	token := os.Getenv("KIBELA_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("KIBELA_TOKEN environment variable is not set")
	}

	team := os.Getenv("KIBELA_TEAM")
	if team == "" {
		return nil, fmt.Errorf("KIBELA_TEAM environment variable is not set")
	}

	return &Config{
		Token: token,
		Team:  team,
	}, nil
}

// Endpoint returns the GraphQL endpoint URL for the team.
func (c *Config) Endpoint() string {
	return fmt.Sprintf("https://%s.kibe.la/api/v1", c.Team)
}
