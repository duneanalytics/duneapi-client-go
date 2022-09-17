package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/dune"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config:", err)
		os.Exit(1)
	}

	// Submitting query for execution
	client := dune.NewDuneClient(cfg.Env)
	execution, err := client.RunQuery(cfg.QueryID, cfg.QueryParameters)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to run query:", err)
		os.Exit(1)
	}

	result, err := execution.WaitGetResults(time.Duration(cfg.PollInterval)*time.Second, cfg.MaxRetries)
	out, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode result as json:", err)
	}

	fmt.Println(string(out))
}
