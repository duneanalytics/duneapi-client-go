package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type StatusResponse struct {
	ExecutionID        string          `json:"execution_id,omitempty"`
	QueryID            int             `json:"query_id,omitempty"`
	State              string          `json:"state,omitempty"`
	SubmittedAt        time.Time       `json:"submitted_at,omitempty"`
	ExecutionStartedAt *time.Time      `json:"execution_started_at,omitempty"`
	ExecutionEndedAt   *time.Time      `json:"execution_ended_at,omitempty"`
	CancelledAt        *time.Time      `json:"cancelled_at,omitempty"`
	Error              *ExecutionError `json:"error,omitempty"`
	ResultMetadata     *ResultMetadata `json:"result_metadata,omitempty"`
}

func (s StatusResponse) HasError() error {
	if s.ExecutionID == "" {
		return errors.New("missing executionID")
	}
	if !strings.HasPrefix(s.State, "QUERY_STATE_") {
		return fmt.Errorf("bad state: %v", s.State)
	}

	if s.State == "QUERY_STATE_COMPLETED" {
		if s.ResultMetadata == nil {
			return errors.New("missing results metadata")
		}
		if s.ExecutionEndedAt == nil {
			return errors.New("missing execution endedAt")
		}
	} else {
		if s.ResultMetadata != nil {
			return fmt.Errorf("cannot have results metadata if state: %v", s.State)
		}
	}

	if s.State == "QUERY_STATE_CANCELLED" {
		if s.CancelledAt == nil {
			return errors.New("missing cancelled at")
		}
	} else {
		if s.CancelledAt != nil {
			return errors.New("field CancelledAt shouldn't be present")
		}
	}
	return nil
}
