package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateMatviewName returns a unique, regex-valid matview name (result_ prefix, lowercase
// alphanumerics/underscores, no trailing underscore).
func generateMatviewName() string {
	return fmt.Sprintf("result_sdk_e2e_%d", time.Now().UnixNano())
}

// TestMaterializedViewLifecycle exercises the full CRUD lifecycle against the real API:
// create a source query -> upsert a matview -> get -> refresh -> delete.
//
// Requirements to run: a DUNE_API_KEY whose plan permits creating materialized views. The
// "small" tier maps to the free engine and is_private=false avoids the private-matview plan gate,
// to maximize compatibility. The matview triggers a real execution and consumes credits.
func TestMaterializedViewLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	// A materialized view needs a saved, non-temporary, non-parameterized source query.
	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     fmt.Sprintf("sdk_e2e_matview_src_%d", time.Now().UnixNano()),
		QuerySQL: "SELECT 1 AS x",
	})
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = client.ArchiveQuery(createResp.QueryID) })

	upsertResp, err := client.UpsertMaterializedView(models.UpsertMaterializedViewRequest{
		Name:        generateMatviewName(),
		QueryID:     createResp.QueryID,
		IsPrivate:   false,
		Performance: "small",
	})
	require.NoError(t, err)
	require.NotEmpty(t, upsertResp.SQLID, "upsert should return the fully-qualified matview name")
	assert.NotEmpty(t, upsertResp.ExecutionID, "creating a matview triggers an execution")

	// The upsert response carries the fully-qualified SQL name; use it for the remaining calls.
	fqName := upsertResp.SQLID
	t.Cleanup(func() { _, _ = client.DeleteMaterializedView(fqName) })

	getResp, err := client.GetMaterializedView(fqName)
	require.NoError(t, err)
	assert.Equal(t, fqName, getResp.SQLID)
	assert.Equal(t, createResp.QueryID, getResp.QueryID)
	assert.False(t, getResp.IsPrivate)

	refreshResp, err := client.RefreshMaterializedView(fqName, models.RefreshMaterializedViewRequest{
		Performance: "small",
	})
	require.NoError(t, err)
	assert.Equal(t, fqName, refreshResp.SQLID)
	assert.NotEmpty(t, refreshResp.ExecutionID, "a refresh triggers a new execution")

	deleteResp, err := client.DeleteMaterializedView(fqName)
	require.NoError(t, err)
	assert.Equal(t, "ok", deleteResp.Message)

	// After deletion the matview is gone.
	_, err = client.GetMaterializedView(fqName)
	assert.Error(t, err, "getting a deleted matview should fail")
}

// TestMaterializedViewWithSchedule verifies that a scheduled matview reports its refresh schedule
// on GET. Requires the schedule field on the GET response (duneapi #1029) to be deployed; against
// an older server getResp.Schedule will be nil and this test will fail on the NotNil assertion.
func TestMaterializedViewWithSchedule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	client := setupClient(t)

	createResp, err := client.CreateQuery(models.CreateQueryRequest{
		Name:     fmt.Sprintf("sdk_e2e_matview_sched_src_%d", time.Now().UnixNano()),
		QuerySQL: "SELECT 1 AS x",
	})
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = client.ArchiveQuery(createResp.QueryID) })

	cron := "0 */6 * * *" // every 6 hours; min interval is 15 minutes
	upsertResp, err := client.UpsertMaterializedView(models.UpsertMaterializedViewRequest{
		Name:           generateMatviewName(),
		QueryID:        createResp.QueryID,
		IsPrivate:      false,
		Performance:    "small",
		CronExpression: &cron,
	})
	require.NoError(t, err)
	fqName := upsertResp.SQLID
	t.Cleanup(func() { _, _ = client.DeleteMaterializedView(fqName) })

	getResp, err := client.GetMaterializedView(fqName)
	require.NoError(t, err)
	require.NotNil(t, getResp.Schedule, "GET should return the refresh schedule (requires duneapi #1029)")
	assert.Equal(t, cron, getResp.Schedule.CronExpression)
	assert.Equal(t, "small", getResp.Schedule.Performance)
}
