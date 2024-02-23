package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResultOptions(t *testing.T) {
	require.Equal(t, "limit=32000", ResultOptions{}.ToURLValues().Encode())
}

func TestResultAddPage(t *testing.T) {
	r := ResultsResponse{}
	r.AddPageResult(&ResultsResponse{
		QueryID: 1,
		State:   "state",
		Result: Result{
			Metadata: ResultMetadata{
				ResultSetBytes: 1,
				RowCount:       1,
				DatapointCount: 1,
				TotalRowCount:  2,
			},
			Rows: []map[string]any{
				{"a": 1},
			},
		},
	})
	require.Equal(t, int64(1), r.QueryID)
	require.Equal(t, "state", r.State)
	require.Equal(t, 1, r.Result.Metadata.RowCount)
	require.Equal(t, 2, r.Result.Metadata.TotalRowCount)
	require.Equal(t, 1, len(r.Result.Rows))
}
