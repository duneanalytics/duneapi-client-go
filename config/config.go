package config

import (
	"fmt"
	"os"
)

const DefaultHost = "https://api.dune.com"

type Env struct {
	APIKey string
	Host   string
}

func getenvOrDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}

	return defaultValue
}

func getenvOrError(key string) (string, error) {
	value, found := os.LookupEnv(key)
	if found {
		return value, nil
	}

	return "", fmt.Errorf("environment variable %s must be set", key)
}

// FromEnvVars populates the config from environment variables
func FromEnvVars() (*Env, error) {
	apiKey, err := getenvOrError("DUNE_API_KEY")
	if err != nil {
		return nil, err
	}
	host := getenvOrDefault("DUNE_API_HOST", DefaultHost)

	return &Env{
		APIKey: apiKey,
		Host:   host,
	}, nil
}

// FromAPIKey generates the config from a passed API key. Uses the default Host
func FromAPIKey(apiKey string) *Env {
	return &Env{
		APIKey: apiKey,
		Host:   DefaultHost,
	}
}
