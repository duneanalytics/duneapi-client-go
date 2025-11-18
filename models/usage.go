package models

type UsageRequest struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

type BillingPeriod struct {
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	CreditsUsed     float64 `json:"credits_used"`
	CreditsIncluded int     `json:"credits_included"`
}

type UsageResponse struct {
	PrivateQueries    int             `json:"private_queries"`
	PrivateDashboards int             `json:"private_dashboards"`
	BytesUsed         int64           `json:"bytes_used"`
	BytesAllowed      int64           `json:"bytes_allowed"`
	BillingPeriods    []BillingPeriod `json:"billing_periods"`
}
