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
