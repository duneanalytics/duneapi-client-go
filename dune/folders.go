package dune

import (
	"fmt"
	"net/url"

	"github.com/duneanalytics/duneapi-client-go/models"
)

// ListFolders lists the folders directly under a folder, or under the library root.
// https://docs.dune.com/api-reference/folders/endpoint/list
func (c *duneClient) ListFolders(opts models.ListFoldersOptions) (*models.ListFoldersResponse, error) {
	u := withQuery(fmt.Sprintf(foldersURLTemplate, c.env.Host), opts.ToURLValues())
	return doJSON[models.ListFoldersResponse](c, "GET", u, nil)
}

// GetFolder returns a folder, its subfolders, and the queries and dashboards inside it.
func (c *duneClient) GetFolder(folderID string, opts models.GetFolderOptions) (*models.GetFolderResponse, error) {
	u := withQuery(fmt.Sprintf(folderURLTemplate, c.env.Host, url.PathEscape(folderID)), opts.ToURLValues())
	return doJSON[models.GetFolderResponse](c, "GET", u, nil)
}

// CreateFolder creates a folder. When a folder already exists at that path the error is an
// *APIError with Code folder_already_exists and the existing folder in Details.
func (c *duneClient) CreateFolder(req models.CreateFolderRequest) (*models.Folder, error) {
	return doJSON[models.Folder](c, "POST", fmt.Sprintf(foldersURLTemplate, c.env.Host), req)
}

// UpdateFolder changes the fields that are set in req and keeps the others.
func (c *duneClient) UpdateFolder(folderID string, req models.UpdateFolderRequest) (*models.Folder, error) {
	u := fmt.Sprintf(folderURLTemplate, c.env.Host, url.PathEscape(folderID))
	return doJSON[models.Folder](c, "PATCH", u, req)
}

// MoveFolder moves a folder under another folder, or to the library root.
func (c *duneClient) MoveFolder(folderID string, req models.MoveFolderRequest) (*models.Folder, error) {
	u := fmt.Sprintf(moveFolderURLTemplate, c.env.Host, url.PathEscape(folderID))
	return doJSON[models.Folder](c, "POST", u, req)
}

// DeleteFolder deletes a folder. A folder that is not empty is refused with an *APIError with
// Code folder_not_empty unless recursive is true, which deletes every subfolder and archives
// the queries and dashboards inside.
func (c *duneClient) DeleteFolder(folderID string, recursive bool) (*models.DeleteFolderResponse, error) {
	u := fmt.Sprintf(folderURLTemplate, c.env.Host, url.PathEscape(folderID))
	if recursive {
		u += "?recursive=true"
	}
	return doJSON[models.DeleteFolderResponse](c, "DELETE", u, nil)
}

// MoveContent moves queries and dashboards into a folder, or out of folders. The move is all
// or nothing.
func (c *duneClient) MoveContent(req models.MoveContentRequest) (*models.MoveContentResponse, error) {
	return doJSON[models.MoveContentResponse](c, "POST", fmt.Sprintf(moveContentURLTemplate, c.env.Host), req)
}
