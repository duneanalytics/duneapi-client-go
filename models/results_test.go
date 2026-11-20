package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResultOptions(t *testing.T) {
	require.Equal(t, "limit=32000", ResultOptions{}.ToURLValues().Encode())
}

func TestResultOptionsShaping(t *testing.T) {
	v := ResultOptions{
		Page:    &ResultPageOption{Offset: 10, Limit: 5},
		Columns: []string{"block_time", "amount usd"},
		Filters: `"amount" > 0 AND "token" IN ('ETH', 'USDC')`,
		SortBy:  []string{`"block_time" DESC`, "amount ASC"},
	}.ToURLValues()
	require.Equal(t, "block_time,amount usd", v.Get("columns"))
	require.Equal(t, `"amount" > 0 AND "token" IN ('ETH', 'USDC')`, v.Get("filters"))
	require.Equal(t, `"block_time" DESC,amount ASC`, v.Get("sort_by"))
	require.Equal(t, "10", v.Get("offset"))
	require.Equal(t, "5", v.Get("limit"))
}

func TestResultOptionsSampleOmitsPagination(t *testing.T) {
	v := ResultOptions{
		Page:        &ResultPageOption{Offset: 10, Limit: 5},
		SampleCount: 100,
	}.ToURLValues()
	require.Equal(t, "sample_count=100", v.Encode())
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
