package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

var executeURLTemplate = "https://%s/api/v1/query/%d/execute"

func (c *client) QueryExecute(queryID int, queryParameters map[string]string) (*models.ExecuteResponse, error) {
	executeURL := fmt.Sprintf("%v/api/v1/query/%d/execute", c.urlBase, queryID)
	jsonData, err := json.Marshal(queryParameters)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", executeURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}

	var executeResp models.ExecuteResponse
	decodeBody(resp, &executeResp)
	if err := executeResp.HasError(); err != nil {
		return nil, err
	}

	return &executeResp, nil
}
