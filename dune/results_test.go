package dune

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

func TestQueryResultsV2KeepsShapingOnEveryPage(t *testing.T) {
	var got []url.Values
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/execution/01EXEC/results", r.URL.Path)
		got = append(got, r.URL.Query())
		if len(got) == 1 {
			io.WriteString(w, `{"execution_id":"01EXEC","query_id":1,"state":"QUERY_STATE_COMPLETED",
				"execution_ended_at":"2026-01-01T00:00:00Z","next_offset":1,
				"result":{"rows":[{"amount":3}],"metadata":{"column_names":["amount"],"row_count":1,"total_row_count":2}}}`)
			return
		}
		io.WriteString(w, `{"execution_id":"01EXEC","query_id":1,"state":"QUERY_STATE_COMPLETED",
			"execution_ended_at":"2026-01-01T00:00:00Z",
			"result":{"rows":[{"amount":2}],"metadata":{"column_names":["amount"],"row_count":1,"total_row_count":2}}}`)
	})

	resp, err := client.QueryResultsV2("01EXEC", models.ResultOptions{
		Columns: []string{"amount"},
		Filters: `"amount" > 1`,
		SortBy:  []string{`"amount" DESC`},
	})
	require.NoError(t, err)
	require.Len(t, resp.Result.Rows, 2)
	require.Len(t, got, 2)
	for _, q := range got {
		require.Equal(t, "amount", q.Get("columns"))
		require.Equal(t, `"amount" > 1`, q.Get("filters"))
		require.Equal(t, `"amount" DESC`, q.Get("sort_by"))
	}
	require.Equal(t, "", got[0].Get("offset"))
	require.Equal(t, "1", got[1].Get("offset"))
}

func TestQueryResultsV2SampleIsOneRequest(t *testing.T) {
	var got []url.Values
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.URL.Query())
		// Even if the server advertised another page, a sample must not be paginated.
		nextOffset := `"next_offset":2,`
		if len(got) > 1 {
			nextOffset = ""
		}
		io.WriteString(w, `{"execution_id":"01EXEC","query_id":1,"state":"QUERY_STATE_COMPLETED",
			"execution_ended_at":"2026-01-01T00:00:00Z",`+nextOffset+`
			"result":{"rows":[{"a":1},{"a":2}],"metadata":{"column_names":["a"],"row_count":2,"total_row_count":50}}}`)
	})

	resp, err := client.QueryResultsV2("01EXEC", models.ResultOptions{SampleCount: 2})
	require.NoError(t, err)
	require.Len(t, resp.Result.Rows, 2)
	require.Len(t, got, 1)
	require.Equal(t, url.Values{"sample_count": {"2"}}, got[0])
}

func TestResultsReportsAnUnreadableResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"execution_id":"01EXEC","state":`))
	})

	resp, err := client.QueryResultsV2("01EXEC", models.ResultOptions{})

	require.Nil(t, resp)
	require.ErrorContains(t, err, "failed to parse response")
}

// A later page that fails to decode must not end the paging with the earlier
// pages' rows, which would look like a complete result.
func TestResultsReportsAnUnreadableLaterPage(t *testing.T) {
	requests := 0
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		requests++
		if requests == 1 {
			_, _ = io.WriteString(w, `{"execution_id":"01EXEC","query_id":1,"state":"QUERY_STATE_COMPLETED",
				"execution_ended_at":"2026-01-01T00:00:00Z","next_offset":1,
				"result":{"rows":[{"amount":3}],"metadata":{"column_names":["amount"],"row_count":1,"total_row_count":2}}}`)
			return
		}
		_, _ = io.WriteString(w, `{"execution_id":"01EXEC","result":{"rows":[`)
	})

	resp, err := client.QueryResultsV2("01EXEC", models.ResultOptions{})

	require.Nil(t, resp)
	require.ErrorContains(t, err, "failed to parse response")
	require.Equal(t, 2, requests)
}
