package models

type QueryParameter struct {
	Key         string   `json:"key"`
	Type        string   `json:"type"`
	Value       string   `json:"value"`
	EnumOptions []string `json:"enumOptions,omitempty"`
}

type CreateQueryRequest struct {
	Name        string           `json:"name"`
	QuerySQL    string           `json:"query_sql"`
	Description string           `json:"description,omitempty"`
	IsPrivate   bool             `json:"is_private,omitempty"`
	IsTemp      bool             `json:"is_temp,omitempty"`
	Parameters  []QueryParameter `json:"parameters,omitempty"`
}

type CreateQueryResponse struct {
	QueryID int `json:"query_id"`
}

type GetQueryResponse struct {
	QueryID      int              `json:"query_id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	QuerySQL     string           `json:"query_sql"`
	Owner        string           `json:"owner"`
	IsPrivate    bool             `json:"is_private"`
	IsArchived   bool             `json:"is_archived"`
	IsUnsaved    bool             `json:"is_unsaved"`
	IsTemp       bool             `json:"is_temp"`
	Version      int              `json:"version"`
	QueryEngine  string           `json:"query_engine"`
	Tags         []string         `json:"tags"`
	Parameters   []QueryParameter `json:"parameters"`
}

type UpdateQueryRequest struct {
	Name        *string          `json:"name,omitempty"`
	QuerySQL    *string          `json:"query_sql,omitempty"`
	Description *string          `json:"description,omitempty"`
	IsPrivate   *bool            `json:"is_private,omitempty"`
	IsArchived  *bool            `json:"is_archived,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Parameters  []QueryParameter `json:"parameters,omitempty"`
}

type UpdateQueryResponse struct {
	QueryID int `json:"query_id"`
}
