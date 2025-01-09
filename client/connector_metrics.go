package client

import (
	"encoding/json"

	"github.com/synadia-io/connect/model"
)

type (
	getConnectorMetricsRequest struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id,omitempty"`
		InstanceId   string `json:"instance_id,omitempty"`
	}
)

func (c *client) GetConnectorMetrics(connectorId, deploymentId, instanceId string, opts ...Opt) ([]*model.InstanceMetric, error) {
	req := getConnectorMetricsRequest{
		ConnectorId:  connectorId,
		DeploymentId: deploymentId,
		InstanceId:   instanceId,
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.METRICS"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp []*model.InstanceMetric
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
