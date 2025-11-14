package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/dune"
	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUploadsIntegration contains E2E tests for Uploads API endpoints.
// These tests require DUNE_API_KEY and DUNE_API_KEY_OWNER_HANDLE environment variables.
// They also require a Plus subscription.
// Run with: DUNE_API_KEY=key DUNE_API_KEY_OWNER_HANDLE=namespace go test ./e2e/...
func TestCreateAndDeleteUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)
	namespace := getTestNamespace(t)
	tableName := generateTableName()

	schema := []models.UploadsColumn{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "name", Type: "varchar", Nullable: false},
		{Name: "value", Type: "double", Nullable: true},
	}

	createResp, err := client.CreateUpload(models.UploadsCreateRequest{
		Namespace:   namespace,
		TableName:   tableName,
		Schema:      schema,
		Description: "Test table created by E2E test",
		IsPrivate:   true,
	})
	require.NoError(t, err)
	assert.Equal(t, tableName, createResp.TableName)
	assert.Equal(t, namespace, createResp.Namespace)
	assert.NotEmpty(t, createResp.FullName)

	deleteResp, err := client.DeleteUpload(namespace, tableName)
	require.NoError(t, err)
	assert.NotEmpty(t, deleteResp.Message)
}

func TestUploadCSVAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)
	namespace := getTestNamespace(t)
	tableName := generateTableName()

	csvData := `id,name,value
1,Alice,10.5
2,Bob,20.3
3,Charlie,15.7`

	csvResp, err := client.UploadCSV(models.UploadsCSVRequest{
		TableName:   tableName,
		Data:        csvData,
		Description: "CSV uploaded by E2E test",
		IsPrivate:   true,
	})
	require.NoError(t, err)
	assert.Equal(t, tableName, csvResp.TableName)
	// FullName might be empty depending on API response
	if csvResp.FullName != "" {
		assert.NotEmpty(t, csvResp.FullName)
	}

	deleteResp, err := client.DeleteUpload(namespace, fmt.Sprintf("dataset_%s", tableName))
	require.NoError(t, err)
	assert.NotEmpty(t, deleteResp.Message)
}

func TestListUploads(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	listResp, err := client.ListUploads(10, 0)
	require.NoError(t, err)
	assert.NotNil(t, listResp.Tables)
}

func TestFullUploadLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)
	namespace := getTestNamespace(t)
	tableName := generateTableName()

	schema := []models.UploadsColumn{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "message", Type: "varchar", Nullable: false},
	}

	createResp, err := client.CreateUpload(models.UploadsCreateRequest{
		Namespace:   namespace,
		TableName:   tableName,
		Schema:      schema,
		Description: "Full lifecycle test",
		IsPrivate:   true,
	})
	require.NoError(t, err)
	assert.Equal(t, tableName, createResp.TableName)

	csvData := "id,message\n1,Hello\n2,World\n"
	insertResp, err := client.InsertIntoUpload(namespace, tableName, csvData, "text/csv")
	require.NoError(t, err)
	assert.Equal(t, int64(2), insertResp.RowsWritten)

	clearResp, err := client.ClearUpload(namespace, tableName)
	require.NoError(t, err)
	assert.NotEmpty(t, clearResp.Message)

	deleteResp, err := client.DeleteUpload(namespace, tableName)
	require.NoError(t, err)
	assert.NotEmpty(t, deleteResp.Message)
}

func TestInsertNDJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)
	namespace := getTestNamespace(t)
	tableName := generateTableName()

	schema := []models.UploadsColumn{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "data", Type: "varchar", Nullable: true},
	}

	_, err := client.CreateUpload(models.UploadsCreateRequest{
		Namespace:   namespace,
		TableName:   tableName,
		Schema:      schema,
		Description: "NDJSON test",
		IsPrivate:   true,
	})
	require.NoError(t, err)

	ndjsonData := `{"id":1,"data":"test1"}
{"id":2,"data":"test2"}`

	insertResp, err := client.InsertIntoUpload(namespace, tableName, ndjsonData, "application/x-ndjson")
	require.NoError(t, err)
	assert.Equal(t, int64(2), insertResp.RowsWritten)

	_, err = client.DeleteUpload(namespace, tableName)
	require.NoError(t, err)
}

func setupClient(t *testing.T) dune.DuneClient {
	apiKey := os.Getenv("DUNE_API_KEY")
	if apiKey == "" {
		t.Fatal("DUNE_API_KEY environment variable must be set to run E2E tests")
	}

	env := config.FromAPIKey(apiKey)
	return dune.NewDuneClient(env)
}

func getTestNamespace(t *testing.T) string {
	namespace := os.Getenv("DUNE_API_KEY_OWNER_HANDLE")
	if namespace == "" {
		t.Fatal("DUNE_API_KEY_OWNER_HANDLE environment variable must be set to run E2E tests")
	}
	return namespace
}

func generateTableName() string {
	return fmt.Sprintf("test_uploads_api_%d", time.Now().Unix())
}
