package models

import (
	"fmt"
	"time"
)

type TableColumn struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Nullable    bool           `json:"nullable"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Description string         `json:"description,omitempty"`
}

type TableOwner struct {
	Handle string `json:"handle"`
	Type   string `json:"type"`
}

type TableListElement struct {
	FullName       string        `json:"full_name"`
	IsPrivate      bool          `json:"is_private"`
	TableSizeBytes string        `json:"table_size_bytes"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	PurgedAt       *time.Time    `json:"purged_at,omitempty"`
	Owner          TableOwner    `json:"owner"`
	Columns        []TableColumn `json:"columns"`
}

type TableListResponse struct {
	Tables []TableListElement `json:"tables"`
}

func (t TableListResponse) HasError() error {
	if t.Tables == nil {
		return fmt.Errorf("missing tables array")
	}
	return nil
}

type TableCreateRequest struct {
	Namespace   string        `json:"namespace"`
	TableName   string        `json:"table_name"`
	Schema      []TableColumn `json:"schema"`
	Description string        `json:"description,omitempty"`
	IsPrivate   bool          `json:"is_private,omitempty"`
}

type TableCreateResponse struct {
	Namespace    string `json:"namespace"`
	TableName    string `json:"table_name"`
	FullName     string `json:"full_name"`
	ExampleQuery string `json:"example_query"`
}

func (t TableCreateResponse) HasError() error {
	if t.Namespace == "" {
		return fmt.Errorf("missing namespace")
	}
	if t.TableName == "" {
		return fmt.Errorf("missing table_name")
	}
	if t.FullName == "" {
		return fmt.Errorf("missing full_name")
	}
	return nil
}

type CSVUploadRequest struct {
	TableName   string `json:"table_name"`
	Data        string `json:"data"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
}

type CSVUploadResponse struct {
	Success   bool   `json:"success"`
	TableName string `json:"table_name"`
}

func (c CSVUploadResponse) HasError() error {
	if !c.Success {
		return fmt.Errorf("CSV upload failed")
	}
	if c.TableName == "" {
		return fmt.Errorf("missing table_name")
	}
	return nil
}

type TableInsertResponse struct {
	Name         string `json:"name"`
	RowsWritten  int64  `json:"rows_written"`
	BytesWritten int64  `json:"bytes_written"`
}

func (t TableInsertResponse) HasError() error {
	if t.Name == "" {
		return fmt.Errorf("missing name")
	}
	return nil
}

type TableDeleteResponse struct {
	Message string `json:"message"`
}

func (t TableDeleteResponse) HasError() error {
	if t.Message == "" {
		return fmt.Errorf("missing message")
	}
	return nil
}

type TableClearResponse struct {
	Message string `json:"message"`
}

func (t TableClearResponse) HasError() error {
	if t.Message == "" {
		return fmt.Errorf("missing message")
	}
	return nil
}
