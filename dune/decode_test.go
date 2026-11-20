package dune

import (
	"io"
	"net/http"
	"testing"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

// Every method that decodes a JSON response must report a body it cannot
// parse, instead of returning an empty response that looks like success.
func TestMethodsReportAnUnreadableResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"execution_id":`)
	})

	calls := map[string]func() (any, error){
		"QueryCancel": func() (any, error) { return nil, client.QueryCancel("01EXEC") },
		"QueryExecute": func() (any, error) {
			return client.QueryExecute(models.ExecuteRequest{QueryID: 1})
		},
		"SQLExecute": func() (any, error) {
			return client.SQLExecute(models.ExecuteSQLRequest{SQL: "SELECT 1"})
		},
		"QueryPipelineExecute": func() (any, error) {
			return client.QueryPipelineExecute(models.PipelineExecuteRequest{})
		},
		"PipelineStatus":   func() (any, error) { return client.PipelineStatus("01PIPE") },
		"QueryStatus":      func() (any, error) { return client.QueryStatus("01EXEC") },
		"GetUsage":         func() (any, error) { return client.GetUsage() },
		"ListUploads":      func() (any, error) { return client.ListUploads(10, 0) },
		"CreateUpload":     func() (any, error) { return client.CreateUpload(models.UploadsCreateRequest{}) },
		"UploadCSV":        func() (any, error) { return client.UploadCSV(models.UploadsCSVRequest{}) },
		"DeleteUpload":     func() (any, error) { return client.DeleteUpload("ns", "t") },
		"ClearUpload":      func() (any, error) { return client.ClearUpload("ns", "t") },
		"InsertIntoUpload": func() (any, error) { return client.InsertIntoUpload("ns", "t", "a\n1", "text/csv") },
		"ListTables":       func() (any, error) { return client.ListTables(10, 0) },
		"CreateTable":      func() (any, error) { return client.CreateTable(models.UploadsCreateRequest{}) },
		"UploadCSVDeprecated": func() (any, error) {
			return client.UploadCSVDeprecated(models.UploadsCSVRequest{})
		},
		"DeleteTable": func() (any, error) { return client.DeleteTable("ns", "t") },
		"ClearTable":  func() (any, error) { return client.ClearTable("ns", "t") },
		"InsertTable": func() (any, error) { return client.InsertTable("ns", "t", "a\n1", "text/csv") },
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			resp, err := call()
			require.ErrorContains(t, err, "failed to parse response")
			require.Nil(t, resp)
		})
	}
}
