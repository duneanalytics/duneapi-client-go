package models

type UpsertMaterializedViewRequest struct {
	Name           string  `json:"name"`
	QueryID        int     `json:"query_id"`
	IsPrivate      bool    `json:"is_private"`
	Performance    string  `json:"performance,omitempty"`
	CronExpression *string `json:"cron_expression,omitempty"`
	ExpiresAt      *string `json:"expires_at,omitempty"`
}

type UpsertMaterializedViewResponse struct {
	// SQLID is the fully-qualified SQL name of the materialized view (e.g. dune.my_team.result_x).
	SQLID       string `json:"name"`
	ExecutionID string `json:"execution_id"`
}

type MaterializedViewSchedule struct {
	CronExpression    string  `json:"cron_expression"`
	Performance       string  `json:"performance,omitempty"`
	NextExecutionTime *string `json:"next_execution_time,omitempty"`
	ExpiresAt         *string `json:"expires_at,omitempty"`
}

type GetMaterializedViewResponse struct {
	// ID and the owner fields are only populated for Dune-team API keys.
	ID               string   `json:"id,omitempty"`
	SQLID            string   `json:"sql_id"`
	OwnerUserID      *int     `json:"owner_user_id,omitempty"`
	OwnerTeamID      *int     `json:"owner_team_id,omitempty"`
	QueryID          int      `json:"query_id"`
	IsPrivate        bool     `json:"is_private"`
	LastExecutionIDs []string `json:"last_execution_ids"`
	TableSizeBytes   int64    `json:"table_size_bytes"`
	// Schedule is nil when the materialized view has no scheduled refresh.
	Schedule *MaterializedViewSchedule `json:"schedule,omitempty"`
}

type MaterializedViewListElement struct {
	ID             string `json:"id"`
	SQLID          string `json:"sql_id"`
	QueryID        int    `json:"query_id"`
	IsPrivate      bool   `json:"is_private"`
	TableSizeBytes int64  `json:"table_size_bytes"`
}

type ListMaterializedViewsResponse struct {
	MaterializedViews []*MaterializedViewListElement `json:"materialized_views"`
	// NextOffset is 0 when there are no further pages.
	NextOffset int32 `json:"next_offset"`
}

type RefreshMaterializedViewRequest struct {
	Performance string `json:"performance,omitempty"`
}

type RefreshMaterializedViewResponse struct {
	SQLID       string `json:"sql_id"`
	ExecutionID string `json:"execution_id"`
}

type DeleteMaterializedViewResponse struct {
	Message string `json:"message"`
}
