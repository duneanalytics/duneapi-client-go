package http

import (
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/models"
)

var cancelURLTemplate = "https://%s/api/v1/execution/%s/cancel"

func QueryCancel(env config.Env, executionID string) (*models.CancelResponse, error) {
	cancelURL := fmt.Sprintf(cancelURLTemplate, env.Host, executionID)
	req, err := http.NewRequest("POST", cancelURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var cancelResp models.CancelResponse
	decodeBody(resp, &cancelResp)
	if err := cancelResp.HasError(); err != nil {
		return nil, err
	}

	return &cancelResp, nil
}
