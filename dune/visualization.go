package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (c *duneClient) CreateVisualization(req models.CreateVisualizationRequest) (*models.CreateVisualizationResponse, error) {
	createURL := fmt.Sprintf(createVisualizationURLTemplate, c.env.Host, req.QueryID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var createResp models.CreateVisualizationResponse
	if err := decodeBody(resp, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *duneClient) GetVisualization(visualizationID int) (*models.GetVisualizationResponse, error) {
	getURL := fmt.Sprintf(visualizationURLTemplate, c.env.Host, visualizationID)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var getResp models.GetVisualizationResponse
	if err := decodeBody(resp, &getResp); err != nil {
		return nil, err
	}

	return &getResp, nil
}

func (c *duneClient) UpdateVisualization(visualizationID int, req models.UpdateVisualizationRequest) (*models.UpdateVisualizationResponse, error) {
	updateURL := fmt.Sprintf(visualizationURLTemplate, c.env.Host, visualizationID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var updateResp models.UpdateVisualizationResponse
	if err := decodeBody(resp, &updateResp); err != nil {
		return nil, err
	}

	return &updateResp, nil
}

func (c *duneClient) DeleteVisualization(visualizationID int) (*models.DeleteVisualizationResponse, error) {
	deleteURL := fmt.Sprintf(visualizationURLTemplate, c.env.Host, visualizationID)

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var deleteResp models.DeleteVisualizationResponse
	if err := decodeBody(resp, &deleteResp); err != nil {
		return nil, err
	}

	return &deleteResp, nil
}

func (c *duneClient) ListQueryVisualizations(queryID, limit, offset int) (*models.ListVisualizationsResponse, error) {
	listURL := fmt.Sprintf(listVisualizationsURLTemplate, c.env.Host, queryID)

	params := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	req, err := http.NewRequest("GET", listURL+params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var listResp models.ListVisualizationsResponse
	if err := decodeBody(resp, &listResp); err != nil {
		return nil, err
	}

	return &listResp, nil
}
