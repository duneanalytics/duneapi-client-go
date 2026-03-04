package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateQueryName() string {
	return fmt.Sprintf("sdk_e2e_test_%d", time.Now().UnixNano())
}

func TestCreateAndGetQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)
	queryName := generateQueryName()

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     queryName,
		QuerySQL: "SELECT 1 AS test_value",
	})
	require.NoError(t, err)
	assert.Greater(t, createResp.QueryID, 0)
	t.Cleanup(func() { client.ArchiveQuery(createResp.QueryID) })

	getResp, err := client.GetQuery(createResp.QueryID)
	require.NoError(t, err)
	assert.Equal(t, createResp.QueryID, getResp.QueryID)
	assert.Equal(t, queryName, getResp.Name)
	assert.Equal(t, "SELECT 1 AS test_value", getResp.QuerySQL)
	assert.False(t, getResp.IsArchived)
}

func TestUpdateQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     generateQueryName(),
		QuerySQL: "SELECT 1",
	})
	require.NoError(t, err)
	t.Cleanup(func() { client.ArchiveQuery(createResp.QueryID) })

	updatedName := generateQueryName() + "_updated"
	updatedSQL := "SELECT 2 AS updated_value"
	updateResp, err := client.UpdateQuery(createResp.QueryID, models.UpdateQueryRequest{
		Name:     &updatedName,
		QuerySQL: &updatedSQL,
	})
	require.NoError(t, err)
	assert.Equal(t, createResp.QueryID, updateResp.QueryID)

	getResp, err := client.GetQuery(createResp.QueryID)
	require.NoError(t, err)
	assert.Equal(t, updatedName, getResp.Name)
	assert.Equal(t, updatedSQL, getResp.QuerySQL)
}

func TestArchiveQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     generateQueryName(),
		QuerySQL: "SELECT 1",
	})
	require.NoError(t, err)

	archiveResp, err := client.ArchiveQuery(createResp.QueryID)
	require.NoError(t, err)
	assert.Equal(t, createResp.QueryID, archiveResp.QueryID)

	getResp, err := client.GetQuery(createResp.QueryID)
	require.NoError(t, err)
	assert.True(t, getResp.IsArchived)
}

func TestCreateQueryWithDescription(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:        generateQueryName(),
		QuerySQL:    "SELECT 1",
		Description: "E2E test query with description",
	})
	require.NoError(t, err)
	t.Cleanup(func() { client.ArchiveQuery(createResp.QueryID) })

	getResp, err := client.GetQuery(createResp.QueryID)
	require.NoError(t, err)
	assert.Equal(t, "E2E test query with description", getResp.Description)
}

func TestCreateTempQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	// Create a temp query
	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     generateQueryName(),
		QuerySQL: "SELECT 1 AS test_value",
		IsTemp:   true,
	})
	require.NoError(t, err)
	assert.Greater(t, createResp.QueryID, 0)

	// Get it back and verify it's marked as temp
	getResp, err := client.GetQuery(createResp.QueryID)
	require.NoError(t, err)
	assert.Equal(t, createResp.QueryID, getResp.QueryID)
	assert.True(t, getResp.IsTemp, "query should have is_temp=true")
	assert.True(t, getResp.IsUnsaved, "query should have is_unsaved=true")
	assert.Equal(t, "SELECT 1 AS test_value", getResp.QuerySQL)
}

func TestQueryExecuteAndGetResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     generateQueryName(),
		QuerySQL: "SELECT 1 AS val",
	})
	require.NoError(t, err)
	t.Cleanup(func() { client.ArchiveQuery(createResp.QueryID) })

	execResp, err := client.QueryExecute(models.ExecuteRequest{
		QueryID: createResp.QueryID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, execResp.ExecutionID)
	assert.Contains(t, execResp.State, "QUERY_STATE_")
}

func TestSQLExecute(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	execResp, err := client.SQLExecute(models.ExecuteSQLRequest{
		SQL: "SELECT 1 AS val",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, execResp.ExecutionID)
	assert.Contains(t, execResp.State, "QUERY_STATE_")
}

func TestRunQueryGetRows(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     generateQueryName(),
		QuerySQL: "SELECT 42 AS answer",
	})
	require.NoError(t, err)
	t.Cleanup(func() { client.ArchiveQuery(createResp.QueryID) })

	rows, err := client.RunQueryGetRows(models.ExecuteRequest{
		QueryID: createResp.QueryID,
	})
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, float64(42), rows[0]["answer"])
}
