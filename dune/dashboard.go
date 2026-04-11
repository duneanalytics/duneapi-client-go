package dune

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duneanalytics/duneapi-client-go/models"
)

func (c *duneClient) CreateDashboard(req models.CreateDashboardRequest) (*models.DashboardResponse, error) {
	createURL := fmt.Sprintf(createDashboardURLTemplate, c.env.Host)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var createResp models.DashboardResponse
	if err := decodeBody(resp, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

func (c *duneClient) GetDashboard(dashboardID int) (*models.DashboardResponse, error) {
	getURL := fmt.Sprintf(dashboardURLTemplate, c.env.Host, dashboardID)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var getResp models.DashboardResponse
	if err := decodeBody(resp, &getResp); err != nil {
		return nil, err
	}

	return &getResp, nil
}

func (c *duneClient) GetDashboardBySlug(ownerHandle, slug string) (*models.DashboardResponse, error) {
	getURL := fmt.Sprintf(dashboardBySlugURLTemplate, c.env.Host, ownerHandle, slug)

	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var getResp models.DashboardResponse
	if err := decodeBody(resp, &getResp); err != nil {
		return nil, err
	}

	return &getResp, nil
}

func (c *duneClient) UpdateDashboard(dashboardID int, req models.UpdateDashboardRequest) (*models.DashboardResponse, error) {
	updateURL := fmt.Sprintf(dashboardURLTemplate, c.env.Host, dashboardID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, httpReq)
	if err != nil {
		return nil, err
	}

	var updateResp models.DashboardResponse
	if err := decodeBody(resp, &updateResp); err != nil {
		return nil, err
	}

	return &updateResp, nil
}

func (c *duneClient) ArchiveDashboard(dashboardID int) (*models.ArchiveDashboardResponse, error) {
	archiveURL := fmt.Sprintf(archiveDashboardURLTemplate, c.env.Host, dashboardID)

	req, err := http.NewRequest("POST", archiveURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest(c.env, req)
	if err != nil {
		return nil, err
	}

	var archiveResp models.ArchiveDashboardResponse
	if err := decodeBody(resp, &archiveResp); err != nil {
		return nil, err
	}

	return &archiveResp, nil
}
