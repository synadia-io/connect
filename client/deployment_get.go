package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	getDeploymentRequest struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id"`
	}

	getDeploymentResponse struct {
		Found bool `json:"found"`

		ConnectorId  string                  `json:"connector_id"`
		DeploymentId string                  `json:"deployment_id"`
		Deployment   *model.DeploymentConfig `json:"deployment"`
		Status       model.DeploymentStatus  `json:"status"`
		Allocations  map[string]int          `json:"allocations,omitempty"`
	}
)

func (c *client) GetDeployment(connectorId string, deploymentId string, opts ...Opt) (*model.Deployment, error) {
	req := getDeploymentRequest{
		ConnectorId:  connectorId,
		DeploymentId: deploymentId,
	}

	b, err := c.Request(c.serviceSubject("DEPLOYMENT.GET"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp getDeploymentResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return &model.Deployment{
		ConnectorId:      resp.ConnectorId,
		DeploymentId:     resp.DeploymentId,
		DeploymentConfig: *resp.Deployment,
		Status:           resp.Status,
		Allocations:      resp.Allocations,
	}, nil

}
