package models

import (
	"net/url"
	"strconv"
	"time"
)

// Folder is a folder in a team's or a user's library.
type Folder struct {
	FolderID string `json:"folder_id"`
	Name     string `json:"name"`
	// Path is the full path from the library root, such as /Research/RWA.
	Path string `json:"path"`
	// ParentFolderID is nil for a folder at the library root.
	ParentFolderID *string   `json:"parent_folder_id"`
	Description    *string   `json:"description"`
	Icon           string    `json:"icon"`
	Color          string    `json:"color"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// FolderSummary is a folder as listed under its parent.
type FolderSummary struct {
	FolderID string `json:"folder_id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	// ItemCount counts the queries and dashboards directly inside the folder.
	ItemCount int `json:"item_count"`
}

// FolderItem is a query or a dashboard inside a folder.
type FolderItem struct {
	// Type is "query" or "dashboard".
	Type string `json:"type"`
	// ID is the query ID or the dashboard ID, depending on Type.
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListFoldersOptions are the filters and paging for GET /v1/folders. Zero values are omitted.
type ListFoldersOptions struct {
	// ParentFolderID lists the subfolders of this folder; empty lists the library root.
	ParentFolderID string
	// NameContains is a case-insensitive substring filter on the folder name.
	NameContains string
	// Limit is the maximum number of folders to return (default 50, max 100).
	Limit  int
	Offset int
}

// ToURLValues encodes the set options as query parameters.
func (o ListFoldersOptions) ToURLValues() url.Values {
	v := url.Values{}
	if o.ParentFolderID != "" {
		v.Set("parent_folder_id", o.ParentFolderID)
	}
	if o.NameContains != "" {
		v.Set("name_contains", o.NameContains)
	}
	addPaging(v, o.Limit, o.Offset)
	return v
}

type ListFoldersResponse struct {
	Folders []FolderSummary `json:"folders"`
	// NextOffset is set when more folders exist; pass it back as Offset.
	NextOffset *int `json:"next_offset,omitempty"`
}

// GetFolderOptions pages through the items of GET /v1/folders/{folder_id}.
type GetFolderOptions struct {
	// Limit is the maximum number of items to return (default 50, max 100).
	Limit  int
	Offset int
}

// ToURLValues encodes the set options as query parameters.
func (o GetFolderOptions) ToURLValues() url.Values {
	v := url.Values{}
	addPaging(v, o.Limit, o.Offset)
	return v
}

type GetFolderResponse struct {
	Folder     Folder          `json:"folder"`
	Subfolders []FolderSummary `json:"subfolders"`
	// HasMoreSubfolders is true when there are more than 100 subfolders; page through them
	// with ListFolders and ParentFolderID.
	HasMoreSubfolders bool         `json:"has_more_subfolders"`
	Items             []FolderItem `json:"items"`
	TotalItems        int          `json:"total_items"`
	// NextOffset is set when more items exist; pass it back as Offset.
	NextOffset *int `json:"next_offset,omitempty"`
}

type CreateFolderRequest struct {
	Name string `json:"name"`
	// ParentFolderID creates the folder inside another folder; nil creates it at the library root.
	ParentFolderID *string `json:"parent_folder_id,omitempty"`
	Description    *string `json:"description,omitempty"`
	// Icon defaults to "folder".
	Icon *string `json:"icon,omitempty"`
	// Color is a preset (green, yellow, red, blue-light, gray, orange, blue, pink, mauve, amber)
	// or a hex value such as #E5477D. Defaults to "gray".
	Color *string `json:"color,omitempty"`
}

// UpdateFolderRequest changes the fields that are set and keeps the others.
type UpdateFolderRequest struct {
	Name *string `json:"name,omitempty"`
	// Description set to an empty string clears it.
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// MoveFolderRequest is the body of POST /v1/folders/{folder_id}/move.
type MoveFolderRequest struct {
	// ParentFolderID is the new parent folder; nil moves the folder to the library root.
	// It is always sent, as null when nil.
	ParentFolderID *string `json:"parent_folder_id"`
}

// MoveContentRequest is the body of POST /v1/content/move.
type MoveContentRequest struct {
	// QueryIDs and DashboardIDs hold at most 100 IDs each.
	QueryIDs     []int64 `json:"query_ids,omitempty"`
	DashboardIDs []int64 `json:"dashboard_ids,omitempty"`
	// TargetFolderID is the destination folder; nil takes the content out of folders.
	// It is always sent, as null when nil.
	TargetFolderID *string `json:"target_folder_id"`
}

type MoveContentResponse struct {
	MovedQueryIDs     []int64 `json:"moved_query_ids"`
	MovedDashboardIDs []int64 `json:"moved_dashboard_ids"`
	// TargetFolder is nil when the content moved to the library root.
	TargetFolder *Folder `json:"target_folder"`
}

type DeleteFolderResponse struct {
	FolderID string `json:"folder_id"`
	Path     string `json:"path"`
	// DeletedFolderCount counts the folder itself and every subfolder deleted with it.
	DeletedFolderCount     int `json:"deleted_folder_count"`
	ArchivedQueryCount     int `json:"archived_query_count"`
	ArchivedDashboardCount int `json:"archived_dashboard_count"`
}

// FolderErrorDetails is the details object of a folder endpoint error. Decode it from
// APIError.Details: folder_already_exists sets ExistingFolderID, and folder_not_empty
// sets the counts a recursive delete would remove.
type FolderErrorDetails struct {
	ExistingFolderID *string `json:"existing_folder_id,omitempty"`
	SubfolderCount   *int    `json:"subfolder_count,omitempty"`
	QueryCount       *int    `json:"query_count,omitempty"`
	DashboardCount   *int    `json:"dashboard_count,omitempty"`
}

func addPaging(v url.Values, limit, offset int) {
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		v.Set("offset", strconv.Itoa(offset))
	}
}
