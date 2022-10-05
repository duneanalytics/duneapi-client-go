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

type Params struct {
	Env             Env
	PollInterval    int               `long:"poll-interval" description:"Interval in seconds for polling for results" default:"5"`                            // nolint:lll
	MaxRetries      int               `long:"max-retries" description:"Max number of get errors tolerated before giving up" default:"5"`                      // nolint:lll
	QueryID         int               `required:"false" short:"q" long:"query-id" description:"The ID of the query to execute"`                               // nolint:lll
	QueryParameters map[string]string `required:"false" short:"p" long:"query-parameter" description:"Parameters to pass to the query in a key:value format"` // nolint:lll
	ExecutionID     string            `required:"false" short:"e" long:"execution-id" description:"ID of an existing execution to check status"`              // nolint:lll
}

// Parse parses all the supplied configurations when used as a CLI
func ParseConfig() (Params, error) {
	var config Params
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
