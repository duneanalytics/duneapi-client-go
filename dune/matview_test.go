package dune

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

func TestUpsertMaterializedView(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody models.UpsertMaterializedViewRequest

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.UpsertMaterializedViewResponse{
			SQLID:       "dune.my_team.result_token_summary",
			ExecutionID: "01HZ065",
		})
	})

	cron := "0 */6 * * *"
	resp, err := client.UpsertMaterializedView(models.UpsertMaterializedViewRequest{
		Name:           "result_token_summary",
		QueryID:        42,
		IsPrivate:      true,
		Performance:    "medium",
		CronExpression: &cron,
	})

	require.NoError(t, err)
	require.Equal(t, "POST", gotMethod)
	require.Equal(t, "/api/v1/materialized-views", gotPath)
	require.Equal(t, "result_token_summary", gotBody.Name)
	require.Equal(t, 42, gotBody.QueryID)
	require.True(t, gotBody.IsPrivate)
	require.Equal(t, "medium", gotBody.Performance)
	require.NotNil(t, gotBody.CronExpression)
	require.Equal(t, "0 */6 * * *", *gotBody.CronExpression)
	require.Equal(t, "dune.my_team.result_token_summary", resp.SQLID)
	require.Equal(t, "01HZ065", resp.ExecutionID)
}

// is_private must always be sent (no omitempty) while an unset cron is omitted.
func TestUpsertMaterializedViewBodyOmitsUnsetFields(t *testing.T) {
	var rawBody map[string]any

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &rawBody)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.UpsertMaterializedViewResponse{})
	})

	_, err := client.UpsertMaterializedView(models.UpsertMaterializedViewRequest{
		Name:      "result_token_summary",
		QueryID:   42,
		IsPrivate: false,
	})

	require.NoError(t, err)
	_, hasCron := rawBody["cron_expression"]
	require.False(t, hasCron, "nil cron_expression must be omitted")
	isPrivate, hasPrivate := rawBody["is_private"]
	require.True(t, hasPrivate, "is_private must always be sent")
	require.Equal(t, false, isPrivate)
}

func TestGetMaterializedView(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.GetMaterializedViewResponse{
			SQLID:          "dune.my_team.result_token_summary",
			QueryID:        42,
			IsPrivate:      true,
			TableSizeBytes: 1024,
			Schedule: &models.MaterializedViewSchedule{
				CronExpression: "0 */6 * * *",
				Performance:    "medium",
			},
		})
	})

	resp, err := client.GetMaterializedView("dune.my_team.result_token_summary")

	require.NoError(t, err)
	require.Equal(t, "GET", gotMethod)
	require.Equal(t, "/api/v1/materialized-views/dune.my_team.result_token_summary", gotPath)
	require.Equal(t, int64(42), resp.QueryID)
	require.True(t, resp.IsPrivate)
	require.Equal(t, int64(1024), resp.TableSizeBytes)
	require.NotNil(t, resp.Schedule)
	require.Equal(t, "0 */6 * * *", resp.Schedule.CronExpression)
	require.Equal(t, "medium", resp.Schedule.Performance)
}

// A response from a server that predates schedule support decodes with a nil Schedule.
func TestGetMaterializedViewWithoutSchedule(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.GetMaterializedViewResponse{
			SQLID:   "dune.my_team.result_token_summary",
			QueryID: 42,
		})
	})

	resp, err := client.GetMaterializedView("dune.my_team.result_token_summary")

	require.NoError(t, err)
	require.Nil(t, resp.Schedule)
}

func TestListMaterializedViews(t *testing.T) {
	var gotPath, gotQuery string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ListMaterializedViewsResponse{
			MaterializedViews: []*models.MaterializedViewListElement{
				{SQLID: "dune.my_team.result_a", QueryID: 1},
				{SQLID: "dune.my_team.result_b", QueryID: 2},
			},
			NextOffset: 2,
		})
	})

	resp, err := client.ListMaterializedViews(2, 0)

	require.NoError(t, err)
	require.Equal(t, "/api/v1/materialized-views", gotPath)
	require.Equal(t, "limit=2&offset=0", gotQuery)
	require.Len(t, resp.MaterializedViews, 2)
	require.Equal(t, "dune.my_team.result_a", resp.MaterializedViews[0].SQLID)
	require.Equal(t, int32(2), resp.NextOffset)
}

func TestRefreshMaterializedView(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody models.RefreshMaterializedViewRequest

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.RefreshMaterializedViewResponse{
			SQLID:       "dune.my_team.result_token_summary",
			ExecutionID: "01HZ999",
		})
	})

	resp, err := client.RefreshMaterializedView(
		"dune.my_team.result_token_summary",
		models.RefreshMaterializedViewRequest{Performance: "large"},
	)

	require.NoError(t, err)
	require.Equal(t, "POST", gotMethod)
	require.Equal(t, "/api/v1/materialized-views/dune.my_team.result_token_summary/refresh", gotPath)
	require.Equal(t, "large", gotBody.Performance)
	require.Equal(t, "01HZ999", resp.ExecutionID)
}

func TestDeleteMaterializedView(t *testing.T) {
	var gotMethod, gotPath string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.DeleteMaterializedViewResponse{Message: "ok"})
	})

	resp, err := client.DeleteMaterializedView("dune.my_team.result_token_summary")

	require.NoError(t, err)
	require.Equal(t, "DELETE", gotMethod)
	require.Equal(t, "/api/v1/materialized-views/dune.my_team.result_token_summary", gotPath)
	require.Equal(t, "ok", resp.Message)
}

func TestMaterializedViewErrorResponseIsWrapped(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Materialized view not found"})
	})

	_, err := client.GetMaterializedView("dune.my_team.result_missing")

	require.Error(t, err)
	require.ErrorIs(t, err, ErrorReqUnsuccessful)
	require.Contains(t, err.Error(), "Materialized view not found")
}
