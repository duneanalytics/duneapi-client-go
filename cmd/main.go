package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/dune"
)

func main() {
	queryID := flag.Int("q", 0, "The ID of the query to execute. Conflicts with -e")
	queryParametersStr := flag.String("p", "{}", "Parameters to pass to the query in JSON format")
	executionID := flag.String("e", "", "ID of an existing execution to check status. Conflicts with -q")
	maxRetries := flag.Int("max-retries", 5, "Max number of errors tolerated before giving up")
	pollInterval := flag.Duration("poll-interval", 5*time.Second, "Interval in seconds for polling for results")

	flag.Parse()

	// Guards against providing either both a query and execution ID or neither.
	if (*executionID == "" && *queryID == 0) || (*executionID != "" && *queryID != 0) {
		fmt.Fprintln(os.Stderr, "must provide exactly one of ExecutionID and QueryID")
		os.Exit(1)
	}

	// Load the API key config from the DUNE_API_KEY environment variable
	env, err := config.FromEnvVars()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	client := dune.NewDuneClient(env)
	var execution dune.Execution

	var queryParameters map[string]any
	err = json.Unmarshal([]byte(*queryParametersStr), &queryParameters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse query parameters: %s\n", err.Error())
	}

	if *executionID == "" {
		// Submitting query for new execution
		execution, err = client.RunQuery(*queryID, queryParameters)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to run query:", err)
			os.Exit(1)
		}
	} else {
		// Using existing execution ID
		execution = dune.NewExecution(client, *executionID)
	}

	// An Execution object provides an interface for getting its state, results,
	// cancelling it, or blocking until it completes, which is what WaitGetResults does
	result, err := execution.WaitGetResults(*pollInterval, *maxRetries)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to retrieve results:", err)
	}

	out, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to encode result as json:", err)
	}

	fmt.Println(string(out))
}
