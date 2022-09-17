package config

import (
	flags "github.com/jessevdk/go-flags"
)

type Query struct {
}

type Env struct {
	APIKey string `required:"true" long:"api-key" env:"DUNE_API_KEY" description:"Your Dune API key"`
	Host   string `long:"host" env:"DUNE_API_HOST" description:"Dune API target host" default:"api.dune.com"`
}

type Config struct {
	Env             Env
	QueryID         []int             `required:"true" short:"q" long:"query-id"`
	QueryParameters map[string]string `required:"false" short:"p" long:"query-parameter"`
}

// Parse parses all the supplied configurations when used as a CLI
func ParseConfig() (Config, error) {
	var config Config
	parser := flags.NewParser(&config, flags.Default)
	_, err := parser.Parse()
	return config, err
}

// ParseEnv parses environment variable config when used as a library
func ParseEnv() (Env, error) {
	var env Env
	parser := flags.NewParser(&env, flags.Default)
	_, err := parser.Parse()
	return env, err
}
