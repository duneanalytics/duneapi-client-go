package dune

import (
	"fmt"

	"github.com/duneanalytics/duneapi-client-go/models"
)

// ListSchedules lists the query, dashboard, and materialized view schedules of the API key's
// owner.
func (c *duneClient) ListSchedules(opts models.ListSchedulesOptions) (*models.ListSchedulesResponse, error) {
	var out models.ListSchedulesResponse
	u := withQuery(fmt.Sprintf(schedulesURLTemplate, c.env.Host), opts.ToURLValues())
	return &out, c.doJSON("GET", u, nil, &out)
}

// SetQuerySchedule creates or replaces the schedule of a saved query.
func (c *duneClient) SetQuerySchedule(queryID int, req models.SetScheduleRequest) (*models.SetScheduleResponse, error) {
	var out models.SetScheduleResponse
	return &out, c.doJSON("PUT", fmt.Sprintf(queryScheduleURLTemplate, c.env.Host, queryID), req, &out)
}

// DeleteQuerySchedule deletes the schedule of a query. Deleted is false when it had none.
func (c *duneClient) DeleteQuerySchedule(queryID int) (*models.DeleteScheduleResponse, error) {
	var out models.DeleteScheduleResponse
	return &out, c.doJSON("DELETE", fmt.Sprintf(queryScheduleURLTemplate, c.env.Host, queryID), nil, &out)
}

// SetDashboardSchedule creates or replaces the refresh schedule of a dashboard.
func (c *duneClient) SetDashboardSchedule(
	dashboardID int, req models.SetScheduleRequest,
) (*models.SetScheduleResponse, error) {
	var out models.SetScheduleResponse
	return &out, c.doJSON("PUT", fmt.Sprintf(dashboardScheduleURLTemplate, c.env.Host, dashboardID), req, &out)
}

// DeleteDashboardSchedule deletes the schedule of a dashboard. Deleted is false when it had none.
func (c *duneClient) DeleteDashboardSchedule(dashboardID int) (*models.DeleteScheduleResponse, error) {
	var out models.DeleteScheduleResponse
	return &out, c.doJSON("DELETE", fmt.Sprintf(dashboardScheduleURLTemplate, c.env.Host, dashboardID), nil, &out)
}
