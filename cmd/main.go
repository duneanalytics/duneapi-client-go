package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/http"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config:", err)
		os.Exit(1)
	}

	// Submitting query for execution
	executionResp, err := http.QueryExecute(cfg.Env, cfg.QueryID[0], cfg.QueryParameters)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to submit query for execution:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Submitted query for execution. ExecutionID:", executionResp.ExecutionID)

	// Waiting for execution to finish
	fmt.Fprintln(os.Stderr, "Waiting for execution to finish...")
	_, err = http.WaitForCompletion(cfg.Env, executionResp.ExecutionID)
	if err != nil {
		os.Exit(1)
	}

	// Get results
	results, err := http.QueryResults(cfg.Env, executionResp.ExecutionID)
	out, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode result as json:", err)
	}

	fmt.Println(string(out))
}
