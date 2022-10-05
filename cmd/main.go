package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/dune"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Usage") {
			os.Exit(1)
		}
		return
	}
	if (cfg.ExecutionID == "" && cfg.QueryID == 0) || (cfg.ExecutionID != "" && cfg.QueryID != 0) {
		fmt.Fprintln(os.Stderr, "must provide exactly one of ExecutionID and QueryID")
		os.Exit(1)
	}
	client := dune.NewDuneClient(cfg.Env)
	var execution dune.Execution

	if cfg.ExecutionID == "" {
		// Submitting query for new execution
		execution, err = client.RunQuery(cfg.QueryID, cfg.QueryParameters)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to run query:", err)
			os.Exit(1)
		}
	} else {
		// Using existing execution
		execution = dune.NewExecution(client, cfg.ExecutionID)
	}

	result, err := execution.WaitGetResults(time.Duration(cfg.PollInterval)*time.Second, cfg.MaxRetries)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to retrieve results:", err)
	}

	out, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to encode result as json:", err)
	}

	fmt.Println(string(out))
}
