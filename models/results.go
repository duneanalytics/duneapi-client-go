package models

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"
)

const LimitRows = 32_000

type ResultMetadata struct {
	ColumnNames         []string `json:"column_names,omitempty"`
	ResultSetBytes      int64    `json:"result_set_bytes,omitempty"`
	RowCount            int      `json:"row_count,omitempty"`
	TotalResultSetBytes int64    `json:"total_result_set_bytes,omitempty"`
	TotalRowCount       int      `json:"total_row_count,omitempty"`
	DatapointCount      int      `json:"datapoint_count,omitempty"`
}

type Result struct {
	Metadata ResultMetadata   `json:"metadata,omitempty"`
	Rows     []map[string]any `json:"rows,omitempty"`
}

type ResultsResponse struct {
	QueryID             int64      `json:"query_id"`
	State               string     `json:"state"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	ExpiresAt           time.Time  `json:"expires_at"`
	ExecutionStartedAt  *time.Time `json:"execution_started_at,omitempty"`
	ExecutionEndedAt    *time.Time `json:"execution_ended_at,omitempty"`
	CancelledAt         *time.Time `json:"cancelled_at,omitempty"`
	Error               *any       `json:"error,omitempty"`
	Result              Result     `json:"result,omitempty"`
	NextOffset          *uint64    `json:"next_offset,omitempty"`
	NextURI             *string    `json:"next_uri,omitempty"`
	IsExecutionFinished bool       `json:"is_execution_finished,omitempty"`
}

func (r ResultsResponse) HasError() error {
	if !strings.HasPrefix(r.State, "QUERY_STATE_") {
		return fmt.Errorf("bad state: %v", r.State)
	}

	if r.State == "QUERY_STATE_COMPLETED" {
		if r.ExecutionEndedAt == nil {
			return errors.New("missing execution endedAt")
		}
		if len(r.Result.Rows) != r.Result.Metadata.RowCount {
			return fmt.Errorf("missmatch row count: len(rows): %v, TotalRowCount: %v",
				len(r.Result.Rows),
				r.Result.Metadata.TotalRowCount,
			)
		}
	} else {
		if r.Result.Rows != nil {
			return fmt.Errorf("cannot have result if state: %v", r.State)
		}
	}

	if r.State == "QUERY_STATE_CANCELLED" {
		if r.CancelledAt == nil {
			return errors.New("missing cancelled at")
		}
	} else {
		if r.CancelledAt != nil {
			return errors.New("field CancelledAt shouldn't be present")
		}
	}
	return nil
}

func (r ResultsResponse) IsEmpty() bool {
	return r.State == "" && r.QueryID == 0 && r.SubmittedAt.Equal(time.Time{})
}

func (r *ResultsResponse) AddPageResult(pageResp *ResultsResponse) {
	if r.IsEmpty() {
		// empty result, copy the first page
		r.QueryID = pageResp.QueryID
		r.State = pageResp.State
		r.SubmittedAt = pageResp.SubmittedAt
		r.ExpiresAt = pageResp.ExpiresAt
		r.ExecutionStartedAt = pageResp.ExecutionStartedAt
		r.ExecutionEndedAt = pageResp.ExecutionEndedAt
		r.CancelledAt = pageResp.CancelledAt
		r.Error = pageResp.Error
		r.NextOffset = pageResp.NextOffset
		r.NextURI = pageResp.NextURI
		r.IsExecutionFinished = pageResp.IsExecutionFinished
		// re-use full result from first page
		r.Result = pageResp.Result
	} else {
		// append rows and the incremental metadata fields
		r.Result.Rows = slices.Concat(r.Result.Rows, pageResp.Result.Rows)
		r.Result.Metadata.ResultSetBytes += pageResp.Result.Metadata.ResultSetBytes
		r.Result.Metadata.RowCount += pageResp.Result.Metadata.RowCount
		r.Result.Metadata.DatapointCount += pageResp.Result.Metadata.DatapointCount
		r.IsExecutionFinished = pageResp.IsExecutionFinished
		r.NextOffset = pageResp.NextOffset
	}
}

// ResultOptions is a struct that contains options for getting a result
type ResultOptions struct {
	// request a specific page of rows
	Page *ResultPageOption
}

func (r ResultOptions) ToURLValues() url.Values {
	v := url.Values{}
	if r.Page != nil {
		if r.Page.Offset > 0 {
			v.Add("offset", fmt.Sprintf("%d", r.Page.Offset))
		}
		limit := r.Page.Limit
		if limit == 0 {
			limit = LimitRows
		}
		v.Add("limit", fmt.Sprintf("%d", limit))
	} else {
		// always paginate the requests
		v.Add("limit", fmt.Sprintf("%d", LimitRows))
	}

	return v
}

// To paginate a large result set
type ResultPageOption struct {
	// we can have more than 2^32 rows, so we need to use int64 for the offset
	Offset uint64
	// assume server can't return more than 2^32 rows
	Limit uint32
}
