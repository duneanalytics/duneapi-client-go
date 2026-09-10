package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

// SubmitContracts submits up to 100 contracts for decoding in one request.
// Each item is validated and queued independently; the response carries one result
// per item, matched by index, so one bad ABI does not fail the batch. Submissions are
// attributed to the user who created the API key.
// https://docs.dune.com/api-reference/contracts/endpoint/decode
func (c *duneClient) SubmitContracts(req models.SubmitContractsRequest) (*models.SubmitContractsResponse, error) {
	submitURL := fmt.Sprintf(submitContractsURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", submitURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var submitResp models.SubmitContractsResponse
	if err := decodeBody(resp, &submitResp); err != nil {
		return nil, err
	}
	if err := submitResp.HasError(); err != nil {
		return nil, err
	}

	return &submitResp, nil
}

// ListContractSubmissions lists the contract decoding submissions made by the user
// who created the API key, newest first. Pass NextCursor from a response back as
// Cursor to fetch the next page.
// https://docs.dune.com/api-reference/contracts/endpoint/list
func (c *duneClient) ListContractSubmissions(
	opts models.ListContractSubmissionsOptions,
) (*models.ListContractSubmissionsResponse, error) {
	listURL := fmt.Sprintf(listContractSubmissionsURLTemplate, c.env.Host)
	if query := opts.ToURLValues().Encode(); query != "" {
		listURL += "?" + query
	}

	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var listResp models.ListContractSubmissionsResponse
	if err := decodeBody(resp, &listResp); err != nil {
		return nil, err
	}
	if err := listResp.HasError(); err != nil {
		return nil, err
	}

	return &listResp, nil
}
