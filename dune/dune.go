package dune

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/models"
)

// DuneClient represents all operations available to call externally
type DuneClient interface {
	// RunQueryGetRows submits a query for execution and returns an Execution object
	RunQuery(queryID int, queryParameters map[string]any) (Execution, error)
	// RunQueryGetRows submits a query for execution, blocks until execution is finished, and returns just the result rows
	RunQueryGetRows(queryID int, queryParameters map[string]any) ([]map[string]any, error)

	// QueryCancel cancels the execution of an execution in the pending or executing state
	QueryCancel(executionID string) error
	// QueryExecute submits a query to execute with the provided parameters
	QueryExecute(queryID int, queryParameters map[string]any) (*models.ExecuteResponse, error)
	// QueryStatus returns the current execution status
	QueryStatus(executionID string) (*models.StatusResponse, error)
	// QueryResults returns the results or status of an execution, depending on whether it has completed
	QueryResults(executionID string) (*models.ResultsResponse, error)
	// QueryResultsCSV returns the results of an execution, as CSV text stream if the execution has completed
	QueryResultsCSV(executionID string) (io.Reader, error)
	// QueryResultsByQueryID returns the results of the lastest execution for a given query ID
	QueryResultsByQueryID(queryID string) (*models.ResultsResponse, error)
	// QueryResultsCSVByQueryID returns the results of the lastest execution for a given query ID, as CSV text stream if the execution has completed
	QueryResultsCSVByQueryID(queryID string) (io.Reader, error)
}

type duneClient struct {
	env *config.Env
}

var cancelURLTemplate = "%s/api/v1/execution/%s/cancel"
var executeURLTemplate = "%s/api/v1/query/%d/execute"
var statusURLTemplate = "%s/api/v1/execution/%s/status"
var executionResultsURLTemplate = "%s/api/v1/execution/%s/results"
var executionResultsCSVURLTemplate = "%s/api/v1/execution/%s/results/csv"
var queryResultsURLTemplate = "%s/api/v1/query/%s/results"
var queryResultsCSVURLTemplate = "%s/api/v1/query/%s/results/csv"

var ErrorRetriesExhausted = errors.New("retries have been exhausted")

// NewDuneClient instantiates a new stateless DuneAPI client. Env contains information about the
// API key and target host (which shouldn't be changed, unless you want to run it through a custom proxy).
func NewDuneClient(env *config.Env) *duneClient {
	return &duneClient{
		env: env,
	}
}

func (c *duneClient) RunQuery(queryID int, queryParameters map[string]any) (Execution, error) {
	resp, err := c.QueryExecute(queryID, queryParameters)
	if err != nil {
		return nil, err
	}

	return &execution{
		client: c,
		ID:     resp.ExecutionID,
	}, nil
}

func (c *duneClient) RunQueryGetRows(queryID int, queryParameters map[string]any) ([]map[string]any, error) {
	execution, err := c.RunQuery(queryID, queryParameters)
	if err != nil {
		return nil, err
	}

	pollInterval := 5 * time.Second
	maxRetries := 10
	resp, err := execution.WaitGetResults(pollInterval, maxRetries)
	if err != nil {
		return nil, err
	}

	return resp.Result.Rows, nil
}

func (c *duneClient) QueryCancel(executionID string) error {
	cancelURL := fmt.Sprintf(cancelURLTemplate, c.env.Host, executionID)
	req, err := http.NewRequest("POST", cancelURL, nil)
	if err != nil {
		return err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return err
	}

	var cancelResp models.CancelResponse
	decodeBody(resp, &cancelResp)
	if err := cancelResp.HasError(); err != nil {
		return err
	}

	return nil
}

func (c *duneClient) QueryExecute(queryID int, queryParameters map[string]any) (*models.ExecuteResponse, error) {
	executeURL := fmt.Sprintf(executeURLTemplate, c.env.Host, queryID)
	jsonData, err := json.Marshal(models.ExecuteRequest{
		QueryParameters: queryParameters,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", executeURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var executeResp models.ExecuteResponse
	decodeBody(resp, &executeResp)
	if err := executeResp.HasError(); err != nil {
		return nil, err
	}

	return &executeResp, nil
}

func (c *duneClient) QueryStatus(executionID string) (*models.StatusResponse, error) {
	statusURL := fmt.Sprintf(statusURLTemplate, c.env.Host, executionID)
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var statusResp models.StatusResponse
	decodeBody(resp, &statusResp)
	if err := statusResp.HasError(); err != nil {
		return nil, err
	}

	return &statusResp, nil
}

func (c *duneClient) QueryResults(executionID string) (*models.ResultsResponse, error) {
	resultsURL := fmt.Sprintf(executionResultsURLTemplate, c.env.Host, executionID)
	req, err := http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var resultsResp models.ResultsResponse
	decodeBody(resp, &resultsResp)
	if err := resultsResp.HasError(); err != nil {
		return nil, err
	}

	return &resultsResp, nil
}

func (c *duneClient) QueryResultsCSV(executionID string) (io.Reader, error) {
	resultsURL := fmt.Sprintf(executionResultsCSVURLTemplate, c.env.Host, executionID)
	req, err := http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	// we read whole result into ram here. if there was a paginated API we wouldn't need to
	var buf bytes.Buffer
	defer resp.Body.Close()
	_, err = buf.ReadFrom(resp.Body)
	return &buf, err
}

func (c *duneClient) QueryResultsByQueryID(queryID string) (*models.ResultsResponse, error) {
	resultsURL := fmt.Sprintf(queryResultsURLTemplate, c.env.Host, queryID)
	req, err := http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var resultsResp models.ResultsResponse
	decodeBody(resp, &resultsResp)
	if err := resultsResp.HasError(); err != nil {
		return nil, err
	}

	return &resultsResp, nil
}

func (c *duneClient) QueryResultsCSVByQueryID(queryID string) (io.Reader, error) {
	resultsURL := fmt.Sprintf(queryResultsCSVURLTemplate, c.env.Host, queryID)
	req, err := http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	// we read whole result into ram here. if there was a paginated API we wouldn't need to
	var buf bytes.Buffer
	defer resp.Body.Close()
	_, err = buf.ReadFrom(resp.Body)
	return &buf, err
}
