package dune

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/duneanalytics/duneapi-client-go/models"
	"github.com/stretchr/testify/require"
)

type recorded struct {
	method, path, query string
	body                map[string]any
	hasBody             bool
}

// recordingClient answers every request with respBody and records the request.
func recordingClient(t *testing.T, respBody string) (*duneClient, *recorded) {
	rec := &recorded{}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		rec.method, rec.path, rec.query = r.Method, r.URL.Path, r.URL.RawQuery
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			rec.hasBody = true
			require.NoError(t, json.Unmarshal(raw, &rec.body))
		}
		_, _ = w.Write([]byte(respBody))
	})
	return client, rec
}

func strPtr(s string) *string { return &s }

func TestFolderRequests(t *testing.T) {
	cases := []struct {
		name       string
		call       func(c *duneClient) error
		method     string
		path       string
		query      string
		body       map[string]any
		noBodySent bool
	}{
		{
			name: "list with filters",
			call: func(c *duneClient) error {
				_, err := c.ListFolders(models.ListFoldersOptions{ParentFolderID: "p1", NameContains: "rwa", Limit: 10, Offset: 20})
				return err
			},
			method: "GET", path: "/api/v1/folders",
			query:      "limit=10&name_contains=rwa&offset=20&parent_folder_id=p1",
			noBodySent: true,
		},
		{
			name: "get escapes the id",
			call: func(c *duneClient) error {
				_, err := c.GetFolder("a/b", models.GetFolderOptions{Limit: 5})
				return err
			},
			method: "GET", path: "/api/v1/folders/a/b", query: "limit=5", noBodySent: true,
		},
		{
			name: "create omits unset fields",
			call: func(c *duneClient) error {
				_, err := c.CreateFolder(models.CreateFolderRequest{Name: "RWA", Color: strPtr("blue")})
				return err
			},
			method: "POST", path: "/api/v1/folders",
			body: map[string]any{"name": "RWA", "color": "blue"},
		},
		{
			name: "update sends an empty description to clear it",
			call: func(c *duneClient) error {
				_, err := c.UpdateFolder("f1", models.UpdateFolderRequest{Description: strPtr("")})
				return err
			},
			method: "PATCH", path: "/api/v1/folders/f1",
			body: map[string]any{"description": ""},
		},
		{
			name: "move to root sends null",
			call: func(c *duneClient) error {
				_, err := c.MoveFolder("f1", models.MoveFolderRequest{})
				return err
			},
			method: "POST", path: "/api/v1/folders/f1/move",
			body: map[string]any{"parent_folder_id": nil},
		},
		{
			name: "delete recursive",
			call: func(c *duneClient) error {
				_, err := c.DeleteFolder("f1", true)
				return err
			},
			method: "DELETE", path: "/api/v1/folders/f1", query: "recursive=true", noBodySent: true,
		},
		{
			name: "delete not recursive sends no flag",
			call: func(c *duneClient) error {
				_, err := c.DeleteFolder("f1", false)
				return err
			},
			method: "DELETE", path: "/api/v1/folders/f1", noBodySent: true,
		},
		{
			name: "move content out of folders sends null target",
			call: func(c *duneClient) error {
				_, err := c.MoveContent(models.MoveContentRequest{QueryIDs: []int64{7}})
				return err
			},
			method: "POST", path: "/api/v1/content/move",
			body: map[string]any{"query_ids": []any{float64(7)}, "target_folder_id": nil},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, rec := recordingClient(t, `{}`)
			require.NoError(t, tc.call(client))
			require.Equal(t, tc.method, rec.method)
			require.Equal(t, tc.path, rec.path)
			require.Equal(t, tc.query, rec.query)
			if tc.noBodySent {
				require.False(t, rec.hasBody)
			} else {
				require.Equal(t, tc.body, rec.body)
			}
		})
	}
}

func TestGetFolderDecodesResponse(t *testing.T) {
	client, _ := recordingClient(t, `{"folder":{"folder_id":"f1","path":"/A","parent_folder_id":null},
		"subfolders":[{"folder_id":"f2","path":"/A/B","item_count":3}],
		"items":[{"type":"query","id":4512345,"name":"q","is_private":true}],
		"total_items":1,"next_offset":50}`)
	resp, err := client.GetFolder("f1", models.GetFolderOptions{})
	require.NoError(t, err)
	require.Equal(t, "/A", resp.Folder.Path)
	require.Nil(t, resp.Folder.ParentFolderID)
	require.Equal(t, 3, resp.Subfolders[0].ItemCount)
	require.Equal(t, int64(4512345), resp.Items[0].ID)
	require.True(t, resp.Items[0].IsPrivate)
	require.Equal(t, 50, *resp.NextOffset)
}

func TestFolderErrorKeepsCodeAndDetails(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Folder /Old is not empty","code":"folder_not_empty",
			"details":{"subfolder_count":1,"query_count":3,"dashboard_count":0}}`))
	})

	_, err := client.DeleteFolder("f1", false)

	require.ErrorIs(t, err, ErrorReqUnsuccessful)
	require.EqualError(t, err, "request was not successful [409]: Folder /Old is not empty")
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, http.StatusConflict, apiErr.StatusCode)
	require.Equal(t, "folder_not_empty", apiErr.Code)
	var details models.FolderErrorDetails
	require.NoError(t, json.Unmarshal(apiErr.Details, &details))
	require.Equal(t, 3, *details.QueryCount)
	require.Equal(t, 0, *details.DashboardCount)
	require.Nil(t, details.ExistingFolderID)
}
