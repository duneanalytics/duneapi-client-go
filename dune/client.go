package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

type QueryExecution interface {
	Cancel(context.Context) (*models.CancelResponse, error)
	Status(context.Context) (*models.StatusResponse, error)
	Result(context.Context) (*models.ResultsResponse, error)
}

type client struct {
	apiKey  string
	urlBase string
}

type Execution struct {
	client *client
	ID     string
}

var ErrorReqUnsuccessful = errors.New("Request was not successful")

func decodeBody(resp *http.Response, dest interface{}) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(dest)
	if err != nil {
		return fmt.Errorf("Failed to parse response: %w", err)
	}
	return nil

}
func (c *client) Request(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-DUNE-API-KEY", c.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %w", err)
	}

	if resp.StatusCode != 200 {
		var errMsg []byte
		defer resp.Body.Close()
		_, err := resp.Body.Read(errMsg)
		if err != nil {
			return nil, fmt.Errorf("Failed to read response body: %w", err)
		}
		return resp, fmt.Errorf("%w [%d]: %s", ErrorReqUnsuccessful, resp.StatusCode, errMsg)
	}

	return resp, nil
}
