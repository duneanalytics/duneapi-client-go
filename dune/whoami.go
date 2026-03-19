package dune

import (
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (c *duneClient) WhoAmI() (*models.WhoAmIResponse, error) {
	reqURL := fmt.Sprintf(whoamiURLTemplate, c.env.Host)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env.APIKey, req)
	if err != nil {
		return nil, err
	}

	var whoamiResp models.WhoAmIResponse
	if err := decodeBody(resp, &whoamiResp); err != nil {
		return nil, err
	}

	return &whoamiResp, nil
}
