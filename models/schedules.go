package models

import (
	"net/url"
	"time"
)

// Schedule runs a saved query, refreshes a dashboard, or refreshes a materialized view on a
// cron schedule. A query or a dashboard has at most one schedule.
type Schedule struct {
	// Type is "query", "dashboard", or "materialized_view". Materialized view refreshes are
	// managed with the materialized view methods.
	Type string `json:"type"`
	// QueryID is set for query and materialized view schedules.
	QueryID *int64 `json:"query_id,omitempty"`
	// DashboardID is set for dashboard schedules.
	DashboardID *int64 `json:"dashboard_id,omitempty"`
	// CronExpression is a 5-part cron expression in UTC.
	CronExpression string `json:"cron_expression"`
	// Performance is the engine tier each run uses: small, medium, or large.
	Performance       string     `json:"performance"`
	NextExecutionTime *time.Time `json:"next_execution_time,omitempty"`
	// ExpiresAt is nil when the schedule does not expire.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ListSchedulesOptions pages through GET /v1/schedules. Zero values are omitted.
type ListSchedulesOptions struct {
	// Limit is the maximum number of schedules to return (default 20, max 100).
	Limit  int
	Offset int
}

// ToURLValues encodes the set options as query parameters.
func (o ListSchedulesOptions) ToURLValues() url.Values {
	v := url.Values{}
	addPaging(v, o.Limit, o.Offset)
	return v
}

type ListSchedulesResponse struct {
	Schedules []Schedule `json:"schedules"`
	// Total is the number of schedules across all pages.
	Total int `json:"total"`
}

// SetScheduleRequest replaces all the settings of a query or dashboard schedule.
type SetScheduleRequest struct {
	// CronExpression is a 5-part cron expression in UTC. Supported: every 15, 20 or 30
	// minutes ("*/15 * * * *"), hourly ("0 * * * *"), every 2, 3, 4, 6, 8 or 12 hours
	// ("0 */6 * * *"), daily ("M H * * *") and weekly ("M H * * D"), where M is 0 or 30.
	CronExpression string `json:"cron_expression"`
	// Performance is small, medium, or large, and must be available on your plan.
	// Every run is charged credits at this tier.
	Performance string `json:"performance"`
	// ExpiresAt stops the schedule after this time; nil means no expiration.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type SetScheduleResponse struct {
	Schedule
	// Created is true when the request created the schedule, false when it replaced one.
	Created bool `json:"created"`
}

type DeleteScheduleResponse struct {
	// Deleted is false when there was no schedule.
	Deleted bool `json:"deleted"`
}
