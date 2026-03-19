package models

// WhoAmIResponse represents the response from the /api/whoami endpoint.
// It identifies the user or team associated with the current API key.
type WhoAmIResponse struct {
	CustomerID string `json:"customer_id"` // e.g. "user_123" or "team_456"
	CreatedBy  int    `json:"created_by"`  // integer user ID of the key creator
	APIKeyID   string `json:"api_key_id"`  // ULID of the API key
}
