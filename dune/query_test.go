package dune

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

func newTestClient(handler http.HandlerFunc) *duneClient {
	server := httptest.NewServer(handler)
	return &duneClient{
		env: &config.Env{
			APIKey: "test-api-key",
			Host:   server.URL,
		},
	}
}

func TestCreateQuery(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody models.CreateQueryRequest

	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.CreateQueryResponse{QueryID: 12345})
	})

	resp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     "Test Query",
		QuerySQL: "SELECT 1",
	})

	require.NoError(t, err)
	require.Equal(t, "POST", gotMethod)
	require.Equal(t, "/api/v1/query", gotPath)
	require.Equal(t, "Test Query", gotBody.Name)
	require.Equal(t, "SELECT 1", gotBody.QuerySQL)
	require.Equal(t, 12345, resp.QueryID)
}

func TestGetQuery(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.GetQueryResponse{
			QueryID:  42,
			Name:     "My Query",
			QuerySQL: "SELECT * FROM ethereum.transactions",
		})
	})

	resp, err := client.GetQuery(42)

	require.NoError(t, err)
	require.Equal(t, "GET", gotMethod)
	require.Equal(t, "/api/v1/query/42", gotPath)
	require.Equal(t, 42, resp.QueryID)
	require.Equal(t, "My Query", resp.Name)
	require.Equal(t, "SELECT * FROM ethereum.transactions", resp.QuerySQL)
}

func TestUpdateQuery(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any

	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.UpdateQueryResponse{QueryID: 42})
	})

	newName := "Updated Name"
	resp, err := client.UpdateQuery(42, models.UpdateQueryRequest{
		Name: &newName,
	})

	require.NoError(t, err)
	require.Equal(t, "PATCH", gotMethod)
	require.Equal(t, "/api/v1/query/42", gotPath)
	require.Equal(t, "Updated Name", gotBody["name"])
	require.Equal(t, 42, resp.QueryID)
}

func TestArchiveQuery(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.UpdateQueryResponse{QueryID: 42})
	})

	resp, err := client.ArchiveQuery(42)

	require.NoError(t, err)
	require.Equal(t, "POST", gotMethod)
	require.Equal(t, "/api/v1/query/42/archive", gotPath)
	require.Equal(t, 42, resp.QueryID)
}

func TestCreateQueryError(t *testing.T) {
	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid query"})
	})

	_, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     "Bad Query",
		QuerySQL: "",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid query")
}

func TestGetQueryError(t *testing.T) {
	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "query not found"})
	})

	_, err := client.GetQuery(99999)

	require.Error(t, err)
	require.Contains(t, err.Error(), "query not found")
}

func TestQueryExecuteWithPerformance(t *testing.T) {
	var gotBody map[string]any

	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ExecuteResponse{
			ExecutionID: "01ABCDEFGHIJKLMNOPQRSTUVWX",
			State:       "QUERY_STATE_PENDING",
		})
	})

	resp, err := client.QueryExecute(models.ExecuteRequest{
		QueryID:         123,
		QueryParameters: map[string]any{"key": "value"},
		Performance:     "large",
	})

	require.NoError(t, err)
	require.Equal(t, "01ABCDEFGHIJKLMNOPQRSTUVWX", resp.ExecutionID)
	require.Equal(t, "large", gotBody["performance"])
	require.NotNil(t, gotBody["query_parameters"])
}

func TestSQLExecuteWithQueryParameters(t *testing.T) {
	var gotBody map[string]any

	client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ExecuteResponse{
			ExecutionID: "01ABCDEFGHIJKLMNOPQRSTUVWX",
			State:       "QUERY_STATE_PENDING",
		})
	})

	resp, err := client.SQLExecute(models.ExecuteSQLRequest{
		SQL:             "SELECT 1",
		Performance:     "medium",
		QueryParameters: map[string]any{"wallet": "0xabc"},
	})

	require.NoError(t, err)
	require.Equal(t, "01ABCDEFGHIJKLMNOPQRSTUVWX", resp.ExecutionID)
	require.Equal(t, "SELECT 1", gotBody["sql"])
	require.Equal(t, "medium", gotBody["performance"])
	require.NotNil(t, gotBody["query_parameters"])
}
