package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (c *duneClient) UpsertMaterializedView(
	req models.UpsertMaterializedViewRequest,
) (*models.UpsertMaterializedViewResponse, error) {
	upsertURL := fmt.Sprintf(matviewsURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", upsertURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var upsertResp models.UpsertMaterializedViewResponse
	if err := decodeBody(resp, &upsertResp); err != nil {
		return nil, err
	}

	return &upsertResp, nil
}

func (c *duneClient) GetMaterializedView(name string) (*models.GetMaterializedViewResponse, error) {
	getURL := fmt.Sprintf(matviewURLTemplate, c.env.Host, url.PathEscape(name))

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var getResp models.GetMaterializedViewResponse
	if err := decodeBody(resp, &getResp); err != nil {
		return nil, err
	}

	return &getResp, nil
}

func (c *duneClient) ListMaterializedViews(limit, offset int) (*models.ListMaterializedViewsResponse, error) {
	listURL := fmt.Sprintf(matviewsURLTemplate, c.env.Host)
	params := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	req, err := http.NewRequest("GET", listURL+params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var listResp models.ListMaterializedViewsResponse
	if err := decodeBody(resp, &listResp); err != nil {
		return nil, err
	}

	return &listResp, nil
}

func (c *duneClient) RefreshMaterializedView(
	name string,
	req models.RefreshMaterializedViewRequest,
) (*models.RefreshMaterializedViewResponse, error) {
	refreshURL := fmt.Sprintf(matviewRefreshURLTemplate, c.env.Host, url.PathEscape(name))

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", refreshURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var refreshResp models.RefreshMaterializedViewResponse
	if err := decodeBody(resp, &refreshResp); err != nil {
		return nil, err
	}

	return &refreshResp, nil
}

func (c *duneClient) DeleteMaterializedView(name string) (*models.DeleteMaterializedViewResponse, error) {
	deleteURL := fmt.Sprintf(matviewURLTemplate, c.env.Host, url.PathEscape(name))

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var deleteResp models.DeleteMaterializedViewResponse
	if err := decodeBody(resp, &deleteResp); err != nil {
		return nil, err
	}

	return &deleteResp, nil
}
