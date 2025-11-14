package dune

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/models"
)

// DuneClient represents all operations available to call externally
type DuneClient interface {
	// New APIs to read results in a more flexible way
	// returns the results or status of an execution, depending on whether it has completed
	QueryResultsV2(executionID string, options models.ResultOptions) (*models.ResultsResponse, error)
	// returns the results of a QueryID, depending on whether it has completed
	ResultsByQueryID(queryID string, options models.ResultOptions) (*models.ResultsResponse, error)

	// RunQueryGetRows submits a query for execution and returns an Execution object
	RunQuery(queryID int, queryParameters map[string]any) (Execution, error)
	// RunQueryGetRows submits a query for execution, blocks until execution is finished, and returns just the result rows
	RunQueryGetRows(queryID int, queryParameters map[string]any) ([]map[string]any, error)

	// QueryCancel cancels the execution of an execution in the pending or executing state
	QueryCancel(executionID string) error

	// QueryExecute submits a query to execute with the provided parameters
	QueryExecute(queryID int, queryParameters map[string]any) (*models.ExecuteResponse, error)

	// SQLExecute executes raw SQL with optional performance parameter
	SQLExecute(sql string, performance string) (*models.ExecuteResponse, error)

	// QueryPipelineExecute submits a query pipeline for execution with optional performance parameter
	QueryPipelineExecute(queryID string, performance string) (*models.PipelineExecuteResponse, error)

	// PipelineStatus returns the current pipeline execution status
	PipelineStatus(pipelineExecutionID string) (*models.PipelineStatusResponse, error)

	// RunSQL submits raw SQL for execution and returns an Execution object
	RunSQL(sql string, performance string) (Execution, error)

	// QueryStatus returns the current execution status
	QueryStatus(executionID string) (*models.StatusResponse, error)

	// QueryResults returns the results or status of an execution, depending on whether it has completed
	// DEPRECATED, use QueryResultsV2 instead
	QueryResults(executionID string) (*models.ResultsResponse, error)

	// QueryResultsCSV returns the results of an execution, as CSV text stream if the execution has completed
	QueryResultsCSV(executionID string) (io.Reader, error)

	// QueryResultsByQueryID returns the results of the lastest execution for a given query ID
	// DEPRECATED, use ResultsByQueryID instead
	QueryResultsByQueryID(queryID string) (*models.ResultsResponse, error)

	// QueryResultsCSVByQueryID returns the results of the lastest execution for a given query ID
	// as CSV text stream if the execution has completed
	QueryResultsCSVByQueryID(queryID string) (io.Reader, error)

	// GetUsage returns usage statistics for the current billing period
	GetUsage() (*models.UsageResponse, error)

	// GetUsageForDates returns usage statistics for a specified time range
	GetUsageForDates(startDate, endDate string) (*models.UsageResponse, error)

	// ListDatasets returns a paginated list of datasets with optional filtering
	ListDatasets(limit, offset int, ownerHandle, datasetType string) (*models.ListDatasetsResponse, error)

	// GetDataset returns detailed information about a specific dataset by slug
	GetDataset(slug string) (*models.DatasetResponse, error)

	// ListUploads returns a paginated list of uploaded tables
	ListUploads(limit, offset int) (*models.TableListResponse, error)

	// CreateTable creates an empty table with defined schema
	CreateTable(req models.TableCreateRequest) (*models.TableCreateResponse, error)

	// UploadCSV uploads CSV data to create a new table
	UploadCSV(req models.CSVUploadRequest) (*models.CSVUploadResponse, error)

	// DeleteTable permanently deletes a table and all its data
	DeleteTable(namespace, tableName string) (*models.TableDeleteResponse, error)

	// ClearTable removes all data from a table while preserving schema
	ClearTable(namespace, tableName string) (*models.TableClearResponse, error)

	// InsertTable inserts data into an existing table (CSV or NDJSON format)
	InsertTable(namespace, tableName, data, contentType string) (*models.TableInsertResponse, error)

	// DEPRECATED: Use ListUploads instead. Will be removed March 1, 2026.
	ListTablesDeprecated(limit, offset int) (*models.TableListResponse, error)

	// DEPRECATED: Use CreateTable instead. Will be removed March 1, 2026.
	CreateTableDeprecated(req models.TableCreateRequest) (*models.TableCreateResponse, error)

	// DEPRECATED: Use UploadCSV instead. Will be removed March 1, 2026.
	UploadCSVDeprecated(req models.CSVUploadRequest) (*models.CSVUploadResponse, error)

	// DEPRECATED: Use DeleteTable instead. Will be removed March 1, 2026.
	DeleteTableDeprecated(namespace, tableName string) (*models.TableDeleteResponse, error)

	// DEPRECATED: Use ClearTable instead. Will be removed March 1, 2026.
	ClearTableDeprecated(namespace, tableName string) (*models.TableClearResponse, error)

	// DEPRECATED: Use InsertTable instead. Will be removed March 1, 2026.
	InsertTableDeprecated(namespace, tableName, data, contentType string) (*models.TableInsertResponse, error)
}

type duneClient struct {
	env *config.Env
}

var (
	cancelURLTemplate                = "%s/api/v1/execution/%s/cancel"
	executeURLTemplate               = "%s/api/v1/query/%d/execute"
	sqlExecuteURLTemplate            = "%s/api/v1/sql/execute"
	pipelineExecuteURLTemplate       = "%s/api/v1/query/%s/pipeline/execute"
	pipelineStatusURLTemplate        = "%s/api/v1/pipelines/executions/%s/status"
	statusURLTemplate                = "%s/api/v1/execution/%s/status"
	executionResultsURLTemplate      = "%s/api/v1/execution/%s/results"
	executionResultsCSVURLTemplate   = "%s/api/v1/execution/%s/results/csv"
	queryResultsURLTemplate          = "%s/api/v1/query/%s/results"
	queryResultsCSVURLTemplate       = "%s/api/v1/query/%s/results/csv"
	usageURLTemplate                 = "%s/api/v1/usage"
	listDatasetsURLTemplate          = "%s/api/v1/datasets"
	getDatasetURLTemplate            = "%s/api/v1/datasets/%s"
	listUploadsURLTemplate           = "%s/api/v1/uploads"
	createTableURLTemplate           = "%s/api/v1/uploads"
	uploadCSVURLTemplate             = "%s/api/v1/uploads/csv"
	deleteTableURLTemplate           = "%s/api/v1/uploads/%s/%s"
	clearTableURLTemplate            = "%s/api/v1/uploads/%s/%s/clear"
	insertTableURLTemplate           = "%s/api/v1/uploads/%s/%s/insert"
	listTablesDeprecatedURLTemplate  = "%s/api/v1/tables"
	createTableDeprecatedURLTemplate = "%s/api/v1/table/create"
	uploadCSVDeprecatedURLTemplate   = "%s/api/v1/table/upload/csv"
	deleteTableDeprecatedURLTemplate = "%s/api/v1/table/%s/%s"
	clearTableDeprecatedURLTemplate  = "%s/api/v1/table/%s/%s/clear"
	insertTableDeprecatedURLTemplate = "%s/api/v1/table/%s/%s/insert"
)

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

func (c *duneClient) RunSQL(sql string, performance string) (Execution, error) {
	resp, err := c.SQLExecute(sql, performance)
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

func (c *duneClient) SQLExecute(sql string, performance string) (*models.ExecuteResponse, error) {
	executeURL := fmt.Sprintf(sqlExecuteURLTemplate, c.env.Host)
	jsonData, err := json.Marshal(models.ExecuteSQLRequest{
		SQL:         sql,
		Performance: performance,
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

func (c *duneClient) QueryPipelineExecute(queryID string, performance string) (*models.PipelineExecuteResponse, error) {
	executeURL := fmt.Sprintf(pipelineExecuteURLTemplate, c.env.Host, queryID)
	jsonData, err := json.Marshal(models.PipelineExecuteRequest{
		Performance: performance,
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

	var pipelineResp models.PipelineExecuteResponse
	decodeBody(resp, &pipelineResp)

	return &pipelineResp, nil
}

func (c *duneClient) PipelineStatus(pipelineExecutionID string) (*models.PipelineStatusResponse, error) {
	statusURL := fmt.Sprintf(pipelineStatusURLTemplate, c.env.Host, pipelineExecutionID)
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var pipelineStatusResp models.PipelineStatusResponse
	decodeBody(resp, &pipelineStatusResp)

	return &pipelineStatusResp, nil
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

func (c *duneClient) getResults(url string, options models.ResultOptions) (*models.ResultsResponse, error) {
	var out models.ResultsResponse

	// track if we have request for a single page
	singlePage := options.Page != nil && (options.Page.Offset > 0 || options.Page.Limit > 0)

	if options.Page == nil {
		options.Page = &models.ResultPageOption{Limit: models.LimitRows}
	}

	for {
		url := fmt.Sprintf("%v?%v", url, options.ToURLValues().Encode())
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := httpRequest(c.env.APIKey, req)
		if err != nil {
			return nil, err
		}

		var pageResp models.ResultsResponse
		decodeBody(resp, &pageResp)
		if err := pageResp.HasError(); err != nil {
			return nil, err
		}
		if singlePage {
			return &pageResp, nil
		}
		out.AddPageResult(&pageResp)

		if pageResp.NextOffset == nil {
			break
		}
		options.Page.Offset = *pageResp.NextOffset
	}

	return &out, nil
}

func (c *duneClient) getResultsCSV(url string) (io.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
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

func (c *duneClient) QueryResultsV2(executionID string, options models.ResultOptions) (*models.ResultsResponse, error) {
	url := fmt.Sprintf(executionResultsURLTemplate, c.env.Host, executionID)
	return c.getResults(url, options)
}

func (c *duneClient) ResultsByQueryID(queryID string, options models.ResultOptions) (*models.ResultsResponse, error) {
	url := fmt.Sprintf(queryResultsURLTemplate, c.env.Host, queryID)
	return c.getResults(url, options)
}

func (c *duneClient) QueryResults(executionID string) (*models.ResultsResponse, error) {
	return c.QueryResultsV2(executionID, models.ResultOptions{})
}

func (c *duneClient) QueryResultsByQueryID(queryID string) (*models.ResultsResponse, error) {
	return c.ResultsByQueryID(queryID, models.ResultOptions{})
}

func (c *duneClient) QueryResultsCSV(executionID string) (io.Reader, error) {
	url := fmt.Sprintf(executionResultsCSVURLTemplate, c.env.Host, executionID)
	return c.getResultsCSV(url)
}

func (c *duneClient) QueryResultsCSVByQueryID(queryID string) (io.Reader, error) {
	url := fmt.Sprintf(queryResultsCSVURLTemplate, c.env.Host, queryID)
	return c.getResultsCSV(url)
}

func (c *duneClient) GetUsage() (*models.UsageResponse, error) {
	return c.getUsage(nil, nil)
}

func (c *duneClient) GetUsageForDates(startDate, endDate string) (*models.UsageResponse, error) {
	return c.getUsage(&startDate, &endDate)
}

func (c *duneClient) getUsage(startDate, endDate *string) (*models.UsageResponse, error) {
	usageURL := fmt.Sprintf(usageURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(models.UsageRequest{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", usageURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var usageResp models.UsageResponse
	decodeBody(resp, &usageResp)

	return &usageResp, nil
}

func (c *duneClient) ListDatasets(
	limit, offset int, ownerHandle, datasetType string,
) (*models.ListDatasetsResponse, error) {
	listURL := fmt.Sprintf(listDatasetsURLTemplate, c.env.Host)

	params := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if ownerHandle != "" {
		params += fmt.Sprintf("&owner_handle=%s", url.QueryEscape(ownerHandle))
	}
	if datasetType != "" {
		params += fmt.Sprintf("&type=%s", url.QueryEscape(datasetType))
	}

	req, err := http.NewRequest("GET", listURL+params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var datasetsResp models.ListDatasetsResponse
	decodeBody(resp, &datasetsResp)
	if err := datasetsResp.HasError(); err != nil {
		return nil, err
	}

	return &datasetsResp, nil
}

func (c *duneClient) GetDataset(slug string) (*models.DatasetResponse, error) {
	getURL := fmt.Sprintf(getDatasetURLTemplate, c.env.Host, slug)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var datasetResp models.DatasetResponse
	decodeBody(resp, &datasetResp)
	if err := datasetResp.HasError(); err != nil {
		return nil, err
	}

	return &datasetResp, nil
}

func (c *duneClient) ListUploads(limit, offset int) (*models.TableListResponse, error) {
	listURL := fmt.Sprintf(listUploadsURLTemplate, c.env.Host)

	params := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	req, err := http.NewRequest("GET", listURL+params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var tablesResp models.TableListResponse
	decodeBody(resp, &tablesResp)
	if err := tablesResp.HasError(); err != nil {
		return nil, err
	}

	return &tablesResp, nil
}

func (c *duneClient) CreateTable(req models.TableCreateRequest) (*models.TableCreateResponse, error) {
	createURL := fmt.Sprintf(createTableURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var createResp models.TableCreateResponse
	decodeBody(resp, &createResp)
	if err := createResp.HasError(); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *duneClient) UploadCSV(req models.CSVUploadRequest) (*models.CSVUploadResponse, error) {
	uploadURL := fmt.Sprintf(uploadCSVURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var uploadResp models.CSVUploadResponse
	decodeBody(resp, &uploadResp)
	if err := uploadResp.HasError(); err != nil {
		return nil, err
	}

	return &uploadResp, nil
}

func (c *duneClient) DeleteTable(namespace, tableName string) (*models.TableDeleteResponse, error) {
	deleteURL := fmt.Sprintf(deleteTableURLTemplate, c.env.Host, url.PathEscape(namespace), url.PathEscape(tableName))

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var deleteResp models.TableDeleteResponse
	decodeBody(resp, &deleteResp)
	if err := deleteResp.HasError(); err != nil {
		return nil, err
	}

	return &deleteResp, nil
}

func (c *duneClient) ClearTable(namespace, tableName string) (*models.TableClearResponse, error) {
	clearURL := fmt.Sprintf(clearTableURLTemplate, c.env.Host, url.PathEscape(namespace), url.PathEscape(tableName))

	req, err := http.NewRequest("POST", clearURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var clearResp models.TableClearResponse
	decodeBody(resp, &clearResp)
	if err := clearResp.HasError(); err != nil {
		return nil, err
	}

	return &clearResp, nil
}

func (c *duneClient) InsertTable(namespace, tableName, data, contentType string) (*models.TableInsertResponse, error) {
	insertURL := fmt.Sprintf(insertTableURLTemplate, c.env.Host, url.PathEscape(namespace), url.PathEscape(tableName))

	req, err := http.NewRequest("POST", insertURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var insertResp models.TableInsertResponse
	decodeBody(resp, &insertResp)
	if err := insertResp.HasError(); err != nil {
		return nil, err
	}

	return &insertResp, nil
}

func (c *duneClient) ListTablesDeprecated(limit, offset int) (*models.TableListResponse, error) {
	listURL := fmt.Sprintf(listTablesDeprecatedURLTemplate, c.env.Host)

	params := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	req, err := http.NewRequest("GET", listURL+params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var tablesResp models.TableListResponse
	decodeBody(resp, &tablesResp)
	if err := tablesResp.HasError(); err != nil {
		return nil, err
	}

	return &tablesResp, nil
}

func (c *duneClient) CreateTableDeprecated(req models.TableCreateRequest) (*models.TableCreateResponse, error) {
	createURL := fmt.Sprintf(createTableDeprecatedURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var createResp models.TableCreateResponse
	decodeBody(resp, &createResp)
	if err := createResp.HasError(); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *duneClient) UploadCSVDeprecated(req models.CSVUploadRequest) (*models.CSVUploadResponse, error) {
	uploadURL := fmt.Sprintf(uploadCSVDeprecatedURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var uploadResp models.CSVUploadResponse
	decodeBody(resp, &uploadResp)
	if err := uploadResp.HasError(); err != nil {
		return nil, err
	}

	return &uploadResp, nil
}

func (c *duneClient) DeleteTableDeprecated(namespace, tableName string) (*models.TableDeleteResponse, error) {
	deleteURL := fmt.Sprintf(deleteTableDeprecatedURLTemplate, c.env.Host, url.PathEscape(namespace), url.PathEscape(tableName))

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var deleteResp models.TableDeleteResponse
	decodeBody(resp, &deleteResp)
	if err := deleteResp.HasError(); err != nil {
		return nil, err
	}

	return &deleteResp, nil
}

func (c *duneClient) ClearTableDeprecated(namespace, tableName string) (*models.TableClearResponse, error) {
	clearURL := fmt.Sprintf(clearTableDeprecatedURLTemplate, c.env.Host, namespace, tableName)

	req, err := http.NewRequest("POST", clearURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var clearResp models.TableClearResponse
	decodeBody(resp, &clearResp)
	if err := clearResp.HasError(); err != nil {
		return nil, err
	}

	return &clearResp, nil
}

func (c *duneClient) InsertTableDeprecated(
	namespace, tableName, data, contentType string,
) (*models.TableInsertResponse, error) {
	insertURL := fmt.Sprintf(insertTableDeprecatedURLTemplate, c.env.Host, url.PathEscape(namespace), url.PathEscape(tableName))

	req, err := http.NewRequest("POST", insertURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var insertResp models.TableInsertResponse
	decodeBody(resp, &insertResp)
	if err := insertResp.HasError(); err != nil {
		return nil, err
	}

	return &insertResp, nil
}
