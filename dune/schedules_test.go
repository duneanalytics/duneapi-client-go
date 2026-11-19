package dune

import (
	"testing"
	"time"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

func TestScheduleAndQueryRequests(t *testing.T) {
	expires := time.Date(2027, 1, 2, 3, 4, 5, 0, time.UTC)
	set := models.SetScheduleRequest{CronExpression: "0 6 * * *", Performance: "medium"}
	cases := []struct {
		name   string
		call   func(c *duneClient) error
		method string
		path   string
		query  string
		body   map[string]any
	}{
		{
			name: "list schedules",
			call: func(c *duneClient) error {
				_, err := c.ListSchedules(models.ListSchedulesOptions{Limit: 100, Offset: 40})
				return err
			},
			method: "GET", path: "/api/v1/schedules", query: "limit=100&offset=40",
		},
		{
			name: "set query schedule with expiry",
			call: func(c *duneClient) error {
				req := set
				req.ExpiresAt = &expires
				_, err := c.SetQuerySchedule(11, req)
				return err
			},
			method: "PUT", path: "/api/v1/queries/11/schedule",
			body: map[string]any{
				"cron_expression": "0 6 * * *", "performance": "medium", "expires_at": "2027-01-02T03:04:05Z",
			},
		},
		{
			name: "set dashboard schedule without expiry",
			call: func(c *duneClient) error {
				_, err := c.SetDashboardSchedule(22, set)
				return err
			},
			method: "PUT", path: "/api/v1/dashboards/22/schedule",
			body: map[string]any{"cron_expression": "0 6 * * *", "performance": "medium"},
		},
		{
			name: "delete query schedule",
			call: func(c *duneClient) error {
				_, err := c.DeleteQuerySchedule(11)
				return err
			},
			method: "DELETE", path: "/api/v1/queries/11/schedule",
		},
		{
			name: "delete dashboard schedule",
			call: func(c *duneClient) error {
				_, err := c.DeleteDashboardSchedule(22)
				return err
			},
			method: "DELETE", path: "/api/v1/dashboards/22/schedule",
		},
		{
			name:   "make private",
			call:   func(c *duneClient) error { _, err := c.MakeQueryPrivate(5); return err },
			method: "POST", path: "/api/v1/query/5/private",
		},
		{
			name:   "make public",
			call:   func(c *duneClient) error { _, err := c.MakeQueryPublic(5); return err },
			method: "POST", path: "/api/v1/query/5/unprivate",
		},
		{
			name:   "unarchive",
			call:   func(c *duneClient) error { _, err := c.UnarchiveQuery(5); return err },
			method: "POST", path: "/api/v1/query/5/unarchive",
		},
		{
			name:   "list queries",
			call:   func(c *duneClient) error { _, err := c.ListQueries(20, 40); return err },
			method: "GET", path: "/api/v1/queries", query: "limit=20&offset=40",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, rec := recordingClient(t, `{}`)
			require.NoError(t, tc.call(client))
			require.Equal(t, tc.method, rec.method)
			require.Equal(t, tc.path, rec.path)
			require.Equal(t, tc.query, rec.query)
			require.Equal(t, tc.body, rec.body)
		})
	}
}

func TestSetScheduleDecodesEmbeddedSchedule(t *testing.T) {
	client, _ := recordingClient(t, `{"type":"query","query_id":11,"cron_expression":"0 6 * * *",
		"performance":"medium","created":true}`)
	resp, err := client.SetQuerySchedule(11, models.SetScheduleRequest{})
	require.NoError(t, err)
	require.True(t, resp.Created)
	require.Equal(t, int64(11), *resp.QueryID)
	require.Equal(t, "0 6 * * *", resp.CronExpression)
	require.Nil(t, resp.ExpiresAt)
}
