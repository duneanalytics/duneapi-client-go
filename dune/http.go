package dune

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrorReqUnsuccessful = errors.New("request was not successful")

type ErrorResponse struct {
	Error string `json:"error"`
}

func decodeBody(resp *http.Response, dest interface{}) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(dest)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}

func httpRequest(apiKey string, req *http.Request) (*http.Response, error) {
	req.Header.Add("X-DUNE-API-KEY", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		var errorResponse ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response body: %w", err)
		}
		return resp, fmt.Errorf("%w [%d]: %s", ErrorReqUnsuccessful, resp.StatusCode, errorResponse.Error)
	}

	return resp, nil
}
