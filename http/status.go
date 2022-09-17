package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/duneanalytics/duneapi-client-go/config"
	"github.com/duneanalytics/duneapi-client-go/models"
)

var statusURLTemplate = "https://%s/api/v1/execution/%s/status"
var ErrorRetriesExhausted = errors.New("Retries have been exhausted")

func QueryStatus(env config.Env, executionID string) (*models.StatusResponse, error) {
	statusURL := fmt.Sprintf(statusURLTemplate, env.Host, executionID)
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpRequest(env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var statusResp models.StatusResponse
	decodeBody(resp, &statusResp)
	if err := statusResp.HasError(); err != nil {
		return nil, err
	}

	return &statusResp, nil
}

func WaitForCompletion(env config.Env, executionID string) (*models.StatusResponse, error) {
	errCount := 0
	for {
		time.Sleep(5 * time.Second)
		statusResp, err := QueryStatus(env, executionID)
		if err != nil {
			errCount += 1
			if errCount > 5 {
				return nil, fmt.Errorf("%w. %s", ErrorRetriesExhausted, err.Error())
			}
			continue
		}

		switch statusResp.State {
		case "QUERY_STATE_COMPLETED", "QUERY_STATE_FAILED", "QUERY_STATE_CANCELLED":
			return statusResp, nil
		}
	}
}
