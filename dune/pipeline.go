package dune

import (
	"github.com/duneanalytics/duneapi-client-go/models"
)

type pipeline struct {
	client DuneClient
	ID     string
}

type Pipeline interface {
	GetStatus() (*models.PipelineStatusResponse, error)
	GetID() string
}

func NewPipeline(client DuneClient, ID string) *pipeline {
	return &pipeline{
		client: client,
		ID:     ID,
	}
}

func (p *pipeline) GetStatus() (*models.PipelineStatusResponse, error) {
	return p.client.PipelineStatus(p.ID)
}

func (p *pipeline) GetID() string {
	return p.ID
}
