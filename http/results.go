package http

import (
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/models"
)

var resultsURLTemplate = "https://%s/api/v1/execution/%s/results"

func QueryResults(env config.Env, executionID string) (*models.ResultsResponse, error) {
	resultsURL := fmt.Sprintf(resultsURLTemplate, env.Host, executionID)
	req, err := http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var resultsResp models.ResultsResponse
	decodeBody(resp, &resultsResp)
	if err := resultsResp.HasError(); err != nil {
		return nil, err
	}

	return &resultsResp, nil
}
