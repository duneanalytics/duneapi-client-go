package models

import "encoding/json"

type SearchDatasetsRequest struct {
	Query          *string  `json:"query,omitempty"`
	Categories     []string `json:"categories,omitempty"`
	Blockchains    []string `json:"blockchains,omitempty"`
	DatasetTypes   []string `json:"dataset_types,omitempty"`
	Schemas        []string `json:"schemas,omitempty"`
	OwnerScope     *string  `json:"owner_scope,omitempty"`
	IncludePrivate *bool    `json:"include_private,omitempty"`
	IncludeSchema  *bool    `json:"include_schema,omitempty"`
	IncludeMetadata *bool   `json:"include_metadata,omitempty"`
	Limit          *int32   `json:"limit,omitempty"`
	Offset         *int32   `json:"offset,omitempty"`
}

type SearchDatasetsResponse struct {
	Total      int32                    `json:"total"`
	Results    []SearchDatasetResult    `json:"results"`
	Pagination SearchDatasetsPagination `json:"pagination"`
}

type SearchDatasetResult struct {
	FullName    string                 `json:"full_name"`
	Category    string                 `json:"category"`
	DatasetType *string                `json:"dataset_type,omitempty"`
	Blockchains []string               `json:"blockchains,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
	OwnerScope  *string                `json:"owner_scope,omitempty"`
	Description *string                `json:"description,omitempty"`
	Schema      json.RawMessage        `json:"schema,omitempty"`
	Metadata    *SearchDatasetMetadata `json:"metadata,omitempty"`
}

type SearchDatasetMetadata struct {
	PageRankScore *float64        `json:"page_rank_score,omitempty"`
	Description   *string         `json:"description,omitempty"`
	AbiType       *string         `json:"abi_type,omitempty"`
	ContractName  *string         `json:"contract_name,omitempty"`
	ProjectName   *string         `json:"project_name,omitempty"`
	SpellType     *string         `json:"spell_type,omitempty"`
	SpellMetadata json.RawMessage `json:"spell_metadata,omitempty"`
}

type SearchDatasetsPagination struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	NextOffset *int32 `json:"next_offset,omitempty"`
	HasMore    bool   `json:"has_more"`
}
