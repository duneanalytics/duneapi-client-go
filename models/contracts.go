package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// ContractSubmissionInput is one contract to submit for decoding. It mirrors the
// form at https://dune.com/contracts/new.
type ContractSubmissionInput struct {
	BlockchainName string `json:"blockchain_name"`
	Address        string `json:"address"`
	ProjectName    string `json:"project_name"`
	ContractName   string `json:"contract_name"`
	// ABI is the contract ABI, either as its JSON array of fragments or as a JSON
	// string containing that array.
	ABI                  json.RawMessage `json:"abi"`
	HasMultipleInstances bool            `json:"has_multiple_instances"`
	IsCreatedByFactory   bool            `json:"is_created_by_factory"`
	IsManualABI          bool            `json:"is_manual_abi"`
	IsProxy              bool            `json:"is_proxy"`
	// SubmissionType is one of new (default), upgrade, rename, delete, other.
	// Upgrades, renames, deletions and other always go to manual review.
	SubmissionType *string `json:"submission_type,omitempty"`
	// ResubmissionReason is required for delete and other.
	ResubmissionReason *string `json:"resubmission_reason,omitempty"`
	// OldProjectName and OldContractName are required for rename and name the
	// contract as it is decoded today.
	OldProjectName  *string `json:"old_project_name,omitempty"`
	OldContractName *string `json:"old_contract_name,omitempty"`
	// IdempotencyKey is a client-chosen key, unique per account, that makes the
	// item safe to retry: a key already used returns the existing submission.
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
}

// SubmitContractsRequest is the body of POST /v1/contracts/decode.
type SubmitContractsRequest struct {
	// Submissions holds between 1 and 100 contracts.
	Submissions []ContractSubmissionInput `json:"submissions"`
}

// ContractSubmissionResult is the per-item outcome of a submission batch, matched
// to the request by Index. On success SubmissionID and Status are set; on failure
// Error describes what to fix.
type ContractSubmissionResult struct {
	Index        int     `json:"index"`
	SubmissionID *string `json:"submission_id,omitempty"`
	Status       *string `json:"status,omitempty"`
	Replayed     bool    `json:"replayed,omitempty"`
	Error        *string `json:"error,omitempty"`
}

// SubmitContractsResponse is the response of POST /v1/contracts/decode.
type SubmitContractsResponse struct {
	Results []ContractSubmissionResult `json:"results"`
}

// HasError reports whether the response is missing its results array.
func (s SubmitContractsResponse) HasError() error {
	if s.Results == nil {
		return fmt.Errorf("missing results array")
	}
	return nil
}

// ListContractSubmissionsOptions are the filters and paging for
// GET /v1/contracts/submissions. Zero values are omitted.
type ListContractSubmissionsOptions struct {
	// Limit is the maximum number of submissions to return (default 50, max 250).
	Limit int
	// Cursor is next_cursor from a previous response, to fetch the next page.
	Cursor         string
	BlockchainName string
	Address        string
	ProjectName    string
	ContractName   string
	// Status is one of pending, approved, rejected, processed, in_progress,
	// cancelled, needs_manual_review.
	Status string
}

// ToURLValues encodes the set options as query parameters.
func (o ListContractSubmissionsOptions) ToURLValues() url.Values {
	values := url.Values{}
	if o.Limit > 0 {
		values.Set("limit", strconv.Itoa(o.Limit))
	}
	for key, value := range map[string]string{
		"cursor":          o.Cursor,
		"blockchain_name": o.BlockchainName,
		"address":         o.Address,
		"project_name":    o.ProjectName,
		"contract_name":   o.ContractName,
		"status":          o.Status,
	} {
		if value != "" {
			values.Set(key, value)
		}
	}
	return values
}

// ContractSubmission is a contract decoding submission as returned by
// GET /v1/contracts/submissions.
type ContractSubmission struct {
	ID             string    `json:"id"`
	BlockchainName string    `json:"blockchain_name"`
	Address        string    `json:"address"`
	ProjectName    string    `json:"project_name"`
	ContractName   string    `json:"contract_name"`
	Status         string    `json:"status"`
	SubmissionType string    `json:"submission_type"`
	Comment        *string   `json:"comment,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IdempotencyKey *string   `json:"idempotency_key,omitempty"`
}

// ListContractSubmissionsResponse is the response of GET /v1/contracts/submissions.
type ListContractSubmissionsResponse struct {
	Submissions []ContractSubmission `json:"submissions"`
	Total       int                  `json:"total"`
	// NextCursor is present when more results exist; pass it back as Cursor.
	NextCursor *string `json:"next_cursor,omitempty"`
}

// HasError reports whether the response is missing its submissions array.
func (l ListContractSubmissionsResponse) HasError() error {
	if l.Submissions == nil {
		return fmt.Errorf("missing submissions array")
	}
	return nil
}
