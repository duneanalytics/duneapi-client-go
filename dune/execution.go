package dune

import (
	"fmt"
	"os"
	"time"

	"github.com/duneanalytics/duneapi-client-go/models"
)

type execution struct {
	client DuneClient
	ID     string
}

type Execution interface {
	// QueryCancel cancels the execution
	Cancel() error
	// QueryResults returns the results or status of the execution, depending on whether it has completed
	GetResults() (*models.ResultsResponse, error)
	// QueryStatus returns the current execution status
	GetStatus() (*models.StatusResponse, error)
	// RunQueryGetResults  blocks until the execution is finished and returns the result
	// maxRetries is used when using the RunQueryToCompletion method, to limit the number of times the method
	// will tolerate API errors before giving up. A value of zero will disable the retry limit.
	// It is recommended to set this to something non-zero, as there is a risk that this will block indefinitely
	// if the Dune API is unreachable or returns an error. The pollInterval determines how long to wait between
	// GetResult requests. It is recommended to set to at least 5 seconds to prevent rate-limiting.
	WaitGetResults(pollInterval time.Duration, maxRetries int) (*models.ResultsResponse, error)
	// GetID returns the execution ID
	GetID() string
}

func NewExecution(client DuneClient, ID string) *execution {
	return &execution{
		client: client,
		ID:     ID,
	}
}

func (e *execution) Cancel() error {
	return e.client.QueryCancel(e.ID)
}

func (e *execution) GetStatus() (*models.StatusResponse, error) {
	return e.client.QueryStatus(e.ID)
}

func (e *execution) GetResults() (*models.ResultsResponse, error) {
	return e.client.QueryResults(e.ID)
}

func (e *execution) WaitGetResults(pollInterval time.Duration, maxRetries int) (*models.ResultsResponse, error) {
	errCount := 0
	for {
		resultsResp, err := e.client.QueryResults(e.ID)
		if err != nil {
			if maxRetries != 0 && errCount > maxRetries {
				return nil, fmt.Errorf("%w. %s", ErrorRetriesExhausted, err.Error())
			}
			fmt.Fprintln(os.Stderr, "failed to retrieve results. Retrying...\n", err)
			errCount += 1
			time.Sleep(pollInterval)
			continue
		}

		switch resultsResp.State {
		case "QUERY_STATE_COMPLETED", "QUERY_STATE_FAILED", "QUERY_STATE_CANCELLED":
			return resultsResp, nil
		default:
			time.Sleep(pollInterval)
			continue
		}
	}
}

func (e *execution) GetID() string {
	return e.ID
}
