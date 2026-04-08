package models

type CreateVisualizationRequest struct {
	QueryID     int            `json:"query_id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description,omitempty"`
	Options     map[string]any `json:"options,omitempty"`
}

type CreateVisualizationResponse struct {
	ID int64 `json:"id"`
}

type GetVisualizationResponse struct {
	ID          int64          `json:"id"`
	QueryID     int64          `json:"query_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Options     map[string]any `json:"options"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type UpdateVisualizationRequest struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Options     map[string]any `json:"options"`
}

type UpdateVisualizationResponse struct {
	ID int64 `json:"id"`
}

type DeleteVisualizationResponse struct {
	OK bool `json:"ok"`
}

type ListVisualizationsResponse struct {
	Results    []VisualizationSummary `json:"results"`
	TotalCount int32                  `json:"total_count"`
}

type VisualizationSummary struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
