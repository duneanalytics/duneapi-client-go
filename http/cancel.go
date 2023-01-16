package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (e *Execution) Cancel(ctx context.Context) (*models.CancelResponse, error) {
	cancelURL := fmt.Sprintf("%v/execution/%v/cancel", e.client.urlBase, e.ID)
	req, err := http.NewRequest("POST", cancelURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Request(req)
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
