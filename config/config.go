package config

import (
	"fmt"
	"os"
)

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

// ParseEnv parses environment variable config when used as a library
func ParseEnv() (*Env, error) {
	apiKey, err := getenvOrError("DUNE_API_KEY")
	if err != nil {
		return nil, err
	}
	host := getenvOrDefault("DUNE_API_HOST", "https://api.dune.com")

	return &Env{
		APIKey: apiKey,
		Host:   host,
	}, nil
}
