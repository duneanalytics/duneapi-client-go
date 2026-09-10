package dune

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

func TestSubmitContracts(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]any

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"results":[
			{"index":0,"submission_id":"sub_1","status":"pending"},
			{"index":1,"submission_id":"sub_0","status":"pending","replayed":true},
			{"index":2,"error":"abi must be valid JSON"}
		]}`)
	})

	upgrade := "upgrade"
	reason := "new implementation"
	key := "k/0"
	resp, err := client.SubmitContracts(models.SubmitContractsRequest{
		Submissions: []models.ContractSubmissionInput{
			{
				BlockchainName: "ethereum",
				Address:        "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984",
				ProjectName:    "uniswap",
				ContractName:   "UniswapToken",
				ABI:            json.RawMessage(`[{"type":"event","name":"Transfer","inputs":[]}]`),
				IsProxy:        true,
				IdempotencyKey: &key,
			},
			{
				BlockchainName:     "ethereum",
				Address:            "0x1",
				ProjectName:        "uniswap",
				ContractName:       "Old",
				ABI:                json.RawMessage(`"[]"`),
				SubmissionType:     &upgrade,
				ResubmissionReason: &reason,
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, "POST", gotMethod)
	require.Equal(t, "/api/v1/contracts/decode", gotPath)

	submissions := gotBody["submissions"].([]any)
	require.Len(t, submissions, 2)
	first := submissions[0].(map[string]any)
	require.Equal(t, []any{map[string]any{"type": "event", "name": "Transfer", "inputs": []any{}}}, first["abi"])
	require.Equal(t, true, first["is_proxy"])
	require.Equal(t, false, first["has_multiple_instances"])
	require.Equal(t, "k/0", first["idempotency_key"])
	require.NotContains(t, first, "submission_type")
	second := submissions[1].(map[string]any)
	require.Equal(t, "[]", second["abi"])
	require.Equal(t, "upgrade", second["submission_type"])
	require.NotContains(t, second, "idempotency_key")

	require.Len(t, resp.Results, 3)
	require.Equal(t, "sub_1", *resp.Results[0].SubmissionID)
	require.Equal(t, "pending", *resp.Results[0].Status)
	require.False(t, resp.Results[0].Replayed)
	require.Nil(t, resp.Results[0].Error)
	require.True(t, resp.Results[1].Replayed)
	require.Equal(t, 2, resp.Results[2].Index)
	require.Nil(t, resp.Results[2].SubmissionID)
	require.Equal(t, "abi must be valid JSON", *resp.Results[2].Error)
}

func TestSubmitContractsReturnsAPIError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, `{"error":"Resubmissions and multi-chain batches require a paid team plan."}`)
	})

	_, err := client.SubmitContracts(models.SubmitContractsRequest{
		Submissions: []models.ContractSubmissionInput{{BlockchainName: "base"}},
	})

	require.ErrorIs(t, err, ErrorReqUnsuccessful)
	require.ErrorContains(t, err, "[403]")
	require.ErrorContains(t, err, "paid team plan")
}

func TestListContractSubmissions(t *testing.T) {
	var gotPath, gotQuery string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"submissions":[{
			"id":"sub_1","blockchain_name":"ethereum","address":"0x1",
			"project_name":"p","contract_name":"c","status":"rejected","submission_type":"new",
			"comment":"There is already a decoded contract at this address.",
			"created_at":"2026-09-10T11:04:18.724658Z","updated_at":"2026-09-10T11:05:18Z",
			"idempotency_key":"k/1"
		}],"total":7,"next_cursor":"abc"}`)
	})

	resp, err := client.ListContractSubmissions(models.ListContractSubmissionsOptions{
		Limit:          1,
		Cursor:         "prev",
		BlockchainName: "ethereum",
		Status:         "rejected",
	})

	require.NoError(t, err)
	require.Equal(t, "/api/v1/contracts/submissions", gotPath)
	require.Equal(t, "blockchain_name=ethereum&cursor=prev&limit=1&status=rejected", gotQuery)
	require.Equal(t, 7, resp.Total)
	require.Equal(t, "abc", *resp.NextCursor)
	require.Len(t, resp.Submissions, 1)
	require.Equal(t, "rejected", resp.Submissions[0].Status)
	require.Equal(t, "There is already a decoded contract at this address.", *resp.Submissions[0].Comment)
	require.Equal(t, "k/1", *resp.Submissions[0].IdempotencyKey)
	require.Equal(t, 2026, resp.Submissions[0].CreatedAt.Year())
}

func TestListContractSubmissionsOmitsUnsetOptions(t *testing.T) {
	var gotQuery string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"submissions":[],"total":0}`)
	})

	resp, err := client.ListContractSubmissions(models.ListContractSubmissionsOptions{})

	require.NoError(t, err)
	require.Equal(t, "", gotQuery)
	require.Nil(t, resp.NextCursor)
	require.Empty(t, resp.Submissions)
}
