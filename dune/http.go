package dune

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	// Message is the error message from the response body.
	Message string
	// Code is a stable, machine-readable error code. Only some endpoints set it.
	Code string
	// Details carries endpoint-specific context; decode it into the endpoint's details type.
	Details json.RawMessage
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s [%d]: %s", ErrorReqUnsuccessful, e.StatusCode, e.Message)
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
		return resp, &APIError{
			StatusCode: resp.StatusCode,
			Message:    errorResponse.Error,
			Code:       errorResponse.Code,
			Details:    errorResponse.Details,
		}
	}

	return resp, nil
}

// doJSON sends body, when it is not nil, as JSON and decodes the response into a T. Like the
// other DuneClient methods, it returns nil with any error.
func doJSON[T any](c *duneClient, method, u string, body any) (*T, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}
	var out T
	if err := decodeBody(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func withQuery(u string, values url.Values) string {
	if q := values.Encode(); q != "" {
		return u + "?" + q
	}
	return u
}
