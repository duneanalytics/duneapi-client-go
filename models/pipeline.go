package models

type PipelineQueryExecutionStatus struct {
	Status      string `json:"status,omitempty"`
	QueryID     int    `json:"query_id,omitempty"`
	ExecutionID string `json:"execution_id,omitempty"`
}

type PipelineNodeExecution struct {
	ID                   int                          `json:"id,omitempty"`
	QueryExecutionStatus PipelineQueryExecutionStatus `json:"query_execution_status,omitempty"`
}

type PipelineStatusResponse struct {
	Status         string                  `json:"status,omitempty"`
	NodeExecutions []PipelineNodeExecution `json:"node_executions,omitempty"`
}

