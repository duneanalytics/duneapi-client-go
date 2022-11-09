package dune

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	// QueryResults returns the results or status of an execution, depending on whether it has completed
	QueryResults(executionID string) (*models.ResultsResponse, error)
	// QueryStatus returns the current execution status
	QueryStatus(executionID string) (*models.StatusResponse, error)
}

type duneClient struct {
	env *config.Env
}

var cancelURLTemplate = "%s/api/v1/execution/%s/cancel"
var executeURLTemplate = "%s/api/v1/query/%d/execute"
var statusURLTemplate = "%s/api/v1/execution/%s/status"
var resultsURLTemplate = "%s/api/v1/execution/%s/results"

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
	resultsURL := fmt.Sprintf(resultsURLTemplate, c.env.Host, executionID)
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
