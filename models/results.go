package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type ResultMetadata struct {
	ColumnNames    []string `json:"column_names,omitempty"`
	ResultSetBytes int64    `json:"result_set_bytes,omitempty"`
	TotalRowCount  int      `json:"total_row_count,omitempty"`
}

type ResultsResponse struct {
	QueryID            int64      `json:"query_id"`
	State              string     `json:"state"`
	SubmittedAt        time.Time  `json:"submitted_at"`
	ExpiresAt          time.Time  `json:"expires_at"`
	ExecutionStartedAt *time.Time `json:"execution_started_at,omitempty"`
	ExecutionEndedAt   *time.Time `json:"execution_ended_at,omitempty"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	Error              *any       `json:"error,omitempty"`
	Result             *struct {
		Metadata ResultMetadata   `json:"metadata,omitempty"`
		Rows     []map[string]any `json:"rows,omitempty"`
	} `json:"result,omitempty"`
}

func (r ResultsResponse) HasError() error {
	if !strings.HasPrefix(r.State, "QUERY_STATE_") {
		return fmt.Errorf("bad state: %v", r.State)
	}

	if r.State == "QUERY_STATE_COMPLETED" {
		if r.Result == nil {
			return errors.New("missing results.result")
		}
		if r.ExecutionEndedAt == nil {
			return errors.New("missing execution endedAt")
		}
		if len(r.Result.Metadata.ColumnNames) == 0 {
			return errors.New("empty column names")
		}
		if r.Result.Metadata.ResultSetBytes == 0 {
			return errors.New("impossible to have ResultSetBytes be 0")
		}
		if len(r.Result.Rows) != r.Result.Metadata.TotalRowCount {
			return fmt.Errorf("missmatch row count: len(rows): %v, TotalRowCount: %v",
				len(r.Result.Rows),
				r.Result.Metadata.TotalRowCount,
			)
		}
	} else {
		if r.Result != nil {
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
