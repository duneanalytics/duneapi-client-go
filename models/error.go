package models

type ExecutionError struct {
	Type     string                 `json:"type"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
