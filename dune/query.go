package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (c *duneClient) CreateQuery(req models.CreateQueryRequest) (*models.CreateQueryResponse, error) {
	createURL := fmt.Sprintf(createQueryURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var createResp models.CreateQueryResponse
	if err := decodeBody(resp, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *duneClient) GetQuery(queryID int) (*models.GetQueryResponse, error) {
	getURL := fmt.Sprintf(queryURLTemplate, c.env.Host, queryID)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var getResp models.GetQueryResponse
	if err := decodeBody(resp, &getResp); err != nil {
		return nil, err
	}

	return &getResp, nil
}

func (c *duneClient) UpdateQuery(queryID int, req models.UpdateQueryRequest) (*models.UpdateQueryResponse, error) {
	updateURL := fmt.Sprintf(queryURLTemplate, c.env.Host, queryID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, httpReq)
	if err != nil {
		return nil, err
	}

	var updateResp models.UpdateQueryResponse
	if err := decodeBody(resp, &updateResp); err != nil {
		return nil, err
	}

	return &updateResp, nil
}

func (c *duneClient) ArchiveQuery(queryID int) (*models.UpdateQueryResponse, error) {
	archiveURL := fmt.Sprintf(archiveQueryURLTemplate, c.env.Host, queryID)

	req, err := http.NewRequest("POST", archiveURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var archiveResp models.UpdateQueryResponse
	if err := decodeBody(resp, &archiveResp); err != nil {
		return nil, err
	}

	return &archiveResp, nil
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
