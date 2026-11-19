package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/duneanalytics/duneapi-client-go/models"
)

// ListFolders lists the folders directly under a folder, or under the library root.
// https://docs.dune.com/api-reference/folders/endpoint/list
func (c *duneClient) ListFolders(opts models.ListFoldersOptions) (*models.ListFoldersResponse, error) {
	var out models.ListFoldersResponse
	u := withQuery(fmt.Sprintf(foldersURLTemplate, c.env.Host), opts.ToURLValues())
	return &out, c.doJSON("GET", u, nil, &out)
}

// GetFolder returns a folder, its subfolders, and the queries and dashboards inside it.
func (c *duneClient) GetFolder(folderID string, opts models.GetFolderOptions) (*models.GetFolderResponse, error) {
	var out models.GetFolderResponse
	u := withQuery(fmt.Sprintf(folderURLTemplate, c.env.Host, url.PathEscape(folderID)), opts.ToURLValues())
	return &out, c.doJSON("GET", u, nil, &out)
}

// CreateFolder creates a folder. When a folder already exists at that path the error is an
// *APIError with Code folder_already_exists and the existing folder in Details.
func (c *duneClient) CreateFolder(req models.CreateFolderRequest) (*models.Folder, error) {
	var out models.Folder
	return &out, c.doJSON("POST", fmt.Sprintf(foldersURLTemplate, c.env.Host), req, &out)
}

// UpdateFolder changes the fields that are set in req and keeps the others.
func (c *duneClient) UpdateFolder(folderID string, req models.UpdateFolderRequest) (*models.Folder, error) {
	var out models.Folder
	u := fmt.Sprintf(folderURLTemplate, c.env.Host, url.PathEscape(folderID))
	return &out, c.doJSON("PATCH", u, req, &out)
}

// MoveFolder moves a folder under another folder, or to the library root.
func (c *duneClient) MoveFolder(folderID string, req models.MoveFolderRequest) (*models.Folder, error) {
	var out models.Folder
	u := fmt.Sprintf(moveFolderURLTemplate, c.env.Host, url.PathEscape(folderID))
	return &out, c.doJSON("POST", u, req, &out)
}

// DeleteFolder deletes a folder. A folder that is not empty is refused with an *APIError with
// Code folder_not_empty unless recursive is true, which deletes every subfolder and archives
// the queries and dashboards inside.
func (c *duneClient) DeleteFolder(folderID string, recursive bool) (*models.DeleteFolderResponse, error) {
	var out models.DeleteFolderResponse
	u := fmt.Sprintf(folderURLTemplate, c.env.Host, url.PathEscape(folderID))
	if recursive {
		u += "?recursive=true"
	}
	return &out, c.doJSON("DELETE", u, nil, &out)
}

// MoveContent moves queries and dashboards into a folder, or out of folders. The move is all
// or nothing.
func (c *duneClient) MoveContent(req models.MoveContentRequest) (*models.MoveContentResponse, error) {
	var out models.MoveContentResponse
	return &out, c.doJSON("POST", fmt.Sprintf(moveContentURLTemplate, c.env.Host), req, &out)
}

// doJSON sends body, when it is not nil, as JSON and decodes the response into dest.
func (c *duneClient) doJSON(method, u string, body any, dest any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := httpRequest(c.env, req)
	if err != nil {
		return err
	}
	return decodeBody(resp, dest)
}

func withQuery(u string, values url.Values) string {
	if q := values.Encode(); q != "" {
		return u + "?" + q
	}
	return u
}
