package models

import (
	"time"
)

type UploadsColumn struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Nullable    bool           `json:"nullable"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Description string         `json:"description,omitempty"`
}

type UploadsOwner struct {
	Handle string `json:"handle"`
	Type   string `json:"type"`
}

type UploadsListElement struct {
	FullName       string          `json:"full_name"`
	IsPrivate      bool            `json:"is_private"`
	TableSizeBytes string          `json:"table_size_bytes"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	PurgedAt       *time.Time      `json:"purged_at,omitempty"`
	Owner          UploadsOwner    `json:"owner"`
	Columns        []UploadsColumn `json:"columns"`
}

type UploadsListResponse struct {
	Tables []UploadsListElement `json:"tables"`
}

type UploadsCreateRequest struct {
	Namespace   string          `json:"namespace"`
	TableName   string          `json:"table_name"`
	Schema      []UploadsColumn `json:"schema"`
	Description string          `json:"description,omitempty"`
	IsPrivate   bool            `json:"is_private,omitempty"`
}

type UploadsCreateResponse struct {
	Namespace      string `json:"namespace"`
	TableName      string `json:"table_name"`
	FullName       string `json:"full_name"`
	ExampleQuery   string `json:"example_query"`
	AlreadyExisted bool   `json:"already_existed"`
	Message        string `json:"message"`
}

type UploadsCSVRequest struct {
	TableName   string `json:"table_name"`
	Data        string `json:"data"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
}

type UploadsCSVResponse struct {
	Success      bool   `json:"success"`
	TableName    string `json:"table_name"`
	FullName     string `json:"full_name"`
	ExampleQuery string `json:"example_query"`
}

type UploadsInsertResponse struct {
	Name         string `json:"name"`
	RowsWritten  int64  `json:"rows_written"`
	BytesWritten int64  `json:"bytes_written"`
}

type UploadsDeleteResponse struct {
	Message string `json:"message"`
}

type UploadsClearResponse struct {
	Message string `json:"message"`
}
