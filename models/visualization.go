package models

type CreateVisualizationRequest struct {
	QueryID     int            `json:"-"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description,omitempty"`
	Options     map[string]any `json:"options,omitempty"`
}

type CreateVisualizationResponse struct {
	ID int64 `json:"id"`
}
