package dune

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/config"
)

var ErrorReqUnsuccessful = errors.New("request was not successful")

type ErrorResponse struct {
	Error string `json:"error"`
	// Code is a stable, machine-readable error code. Only some endpoints, such as the
	// folder endpoints, set it.
	Code string `json:"code,omitempty"`
	// Details carries endpoint-specific context, such as the existing folder for
	// folder_already_exists. Decode it into the endpoint's details type.
	Details json.RawMessage `json:"details,omitempty"`
}

// APIError is returned for every non-2xx response. It wraps ErrorReqUnsuccessful,
// so errors.Is(err, ErrorReqUnsuccessful) keeps working; use errors.As to read the
// status code and the error body.
type APIError struct {
	StatusCode int
	ErrorResponse
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s [%d]: %s", ErrorReqUnsuccessful, e.StatusCode, e.ErrorResponse.Error)
}

func (e *APIError) Unwrap() error {
	return ErrorReqUnsuccessful
}

func decodeBody(resp *http.Response, dest interface{}) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(dest)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}

func httpRequest(env *config.Env, req *http.Request) (*http.Response, error) {
	req.Header.Add("X-DUNE-API-KEY", env.APIKey)
	for key, value := range env.Headers {
		req.Header.Set(key, value)
	}
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
		return resp, &APIError{StatusCode: resp.StatusCode, ErrorResponse: errorResponse}
	}

	return resp, nil
}
