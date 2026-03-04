package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (c *duneClient) ListDatasets(
	limit, offset int, ownerHandle, datasetType string,
) (*models.ListDatasetsResponse, error) {
	listURL := fmt.Sprintf(listDatasetsURLTemplate, c.env.Host)

	params := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if ownerHandle != "" {
		params += fmt.Sprintf("&owner_handle=%s", url.QueryEscape(ownerHandle))
	}
	if datasetType != "" {
		params += fmt.Sprintf("&type=%s", url.QueryEscape(datasetType))
	}

	req, err := http.NewRequest("GET", listURL+params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var datasetsResp models.ListDatasetsResponse
	if err := decodeBody(resp, &datasetsResp); err != nil {
		return nil, err
	}
	if err := datasetsResp.HasError(); err != nil {
		return nil, err
	}

	return &datasetsResp, nil
}

func (c *duneClient) GetDataset(slug string) (*models.DatasetResponse, error) {
	getURL := fmt.Sprintf(getDatasetURLTemplate, c.env.Host, slug)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var datasetResp models.DatasetResponse
	if err := decodeBody(resp, &datasetResp); err != nil {
		return nil, err
	}
	if err := datasetResp.HasError(); err != nil {
		return nil, err
	}

	return &datasetResp, nil
}

func (c *duneClient) SearchDatasets(req models.SearchDatasetsRequest) (*models.SearchDatasetsResponse, error) {
	searchURL := fmt.Sprintf(searchDatasetsURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", searchURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var searchResp models.SearchDatasetsResponse
	if err := decodeBody(resp, &searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}
