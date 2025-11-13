package models

import (
	"fmt"
	"strings"
)

type ExecuteRequest struct {
	QueryParameters map[string]any `json:"query_parameters,omitempty"`
}

type ExecuteSQLRequest struct {
	SQL         string `json:"sql"`
	Performance string `json:"performance,omitempty"`
}

type ExecuteResponse struct {
	ExecutionID string `json:"execution_id,omitempty"`
	State       string `json:"state,omitempty"`
}

func (e ExecuteResponse) HasError() error {
	// 01 is the expected prefix for an ULID string
	if len(e.ExecutionID) != 26 || !strings.HasPrefix(e.ExecutionID, "01") {
		return fmt.Errorf("bad execution id: %v", e.ExecutionID)
	}
	if !strings.HasPrefix(e.State, "QUERY_STATE_") {
		return fmt.Errorf("bad state: %v", e.State)
	}
	return nil
}

type PipelineExecuteRequest struct {
	Performance string `json:"performance,omitempty"`
}

type PipelineExecuteResponse struct {
	PipelineExecutionID string `json:"pipeline_execution_id,omitempty"`
}
