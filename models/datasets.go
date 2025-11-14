package models

import (
	"fmt"
	"time"
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
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Namespace   string          `json:"namespace"`
	TableName   string          `json:"table_name"`
	Type        string          `json:"type"`
	Columns     []DatasetColumn `json:"columns"`
	Owner       DatasetOwner    `json:"owner"`
	IsPrivate   bool            `json:"is_private"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Description string          `json:"description,omitempty"`
}

func (d DatasetResponse) HasError() error {
	if d.Slug == "" {
		return fmt.Errorf("missing dataset slug")
	}
	if d.Namespace == "" {
		return fmt.Errorf("missing namespace")
	}
	if d.TableName == "" {
		return fmt.Errorf("missing table_name")
	}
	if d.Owner.Handle == "" {
		return fmt.Errorf("missing owner handle")
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
