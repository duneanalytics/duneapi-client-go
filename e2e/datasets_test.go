package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatasetsIntegration contains E2E tests for Datasets API endpoints.
// These tests require DUNE_API_KEY and DUNE_API_KEY_OWNER_HANDLE environment variables.
// Run with: DUNE_API_KEY=key DUNE_API_KEY_OWNER_HANDLE=namespace go test ./e2e/...
func TestListDatasets(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	result, err := client.ListDatasets(10, 0, "", "uploaded_table")
	require.NoError(t, err)
	assert.NotNil(t, result.Datasets)
	assert.GreaterOrEqual(t, result.Total, 0)

	if len(result.Datasets) > 0 {
		dataset := result.Datasets[0]
		assert.NotEmpty(t, dataset.FullName)
		assert.NotEmpty(t, dataset.Type)
		assert.NotNil(t, dataset.Owner)
		assert.NotNil(t, dataset.Columns)
	}
}

func TestListDatasetsWithFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	result, err := client.ListDatasets(5, 0, "", "transformation_view")
	require.NoError(t, err)
	assert.NotNil(t, result.Datasets)

	for _, dataset := range result.Datasets {
		assert.Equal(t, "transformation_view", dataset.Type)
	}
}

func TestListDatasetsByOwner(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	result, err := client.ListDatasets(5, 0, "dune", "")
	require.NoError(t, err)
	assert.NotNil(t, result.Datasets)

	for _, dataset := range result.Datasets {
		if dataset.Owner != nil {
			assert.Equal(t, "dune", dataset.Owner.Handle)
		}
	}
}

func TestGetDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	listResult, err := client.ListDatasets(1, 0, "", "uploaded_table")
	require.NoError(t, err)

	if len(listResult.Datasets) == 0 {
		t.Skip("No uploaded tables found to test")
	}

	fullName := listResult.Datasets[0].FullName
	result, err := client.GetDataset(fullName)
	require.NoError(t, err)

	assert.Equal(t, fullName, result.FullName)
	assert.NotEmpty(t, result.Type)
	assert.NotNil(t, result.Owner)
	assert.NotNil(t, result.Columns)
	assert.Greater(t, len(result.Columns), 0)

	column := result.Columns[0]
	assert.NotEmpty(t, column.Name)
	assert.NotEmpty(t, column.Type)
}

func TestGetDatasetWithUploadedTable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	listResult, err := client.ListDatasets(1, 0, "", "uploaded_table")
	require.NoError(t, err)

	if len(listResult.Datasets) > 0 {
		fullName := listResult.Datasets[0].FullName
		result, err := client.GetDataset(fullName)
		require.NoError(t, err)

		assert.Equal(t, fullName, result.FullName)
		assert.Equal(t, "uploaded_table", result.Type)
	}
}
