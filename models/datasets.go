package models

import (
	"fmt"
)

type DatasetColumn struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Nullable    bool           `json:"nullable"`
	Description string         `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type DatasetOwner struct {
	Handle string `json:"handle"`
	Type   string `json:"type"`
}

type DatasetResponse struct {
	Type      string          `json:"type"`
	FullName  string          `json:"full_name"`
	IsPrivate bool            `json:"is_private"`
	Columns   []DatasetColumn `json:"columns"`
	Owner     *DatasetOwner   `json:"owner"`
	Metadata  map[string]any  `json:"metadata,omitempty"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func (d DatasetResponse) HasError() error {
	if d.FullName == "" {
		return fmt.Errorf("missing full_name")
	}
	if d.Type == "" {
		return fmt.Errorf("missing type")
	}
	return nil
}

type ListDatasetsResponse struct {
	Datasets []DatasetResponse `json:"datasets"`
	Total    int               `json:"total"`
}

func (l ListDatasetsResponse) HasError() error {
	if l.Datasets == nil {
		return fmt.Errorf("missing datasets array")
	}
	return nil
}
