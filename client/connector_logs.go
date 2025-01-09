package client

import (
	"encoding/json"

	"github.com/synadia-io/connect/model"
)

type (
	getConnectorLogsRequest struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id,omitempty"`
		InstanceId   string `json:"instance_id,omitempty"`
	}
)

func (c *client) GetLogs(connectorId, deploymentId, instanceId string, opts ...Opt) ([]*model.InstanceLog, error) {
	req := getConnectorLogsRequest{
		ConnectorId:  connectorId,
		DeploymentId: deploymentId,
		InstanceId:   instanceId,
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.LOGS"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp []*model.InstanceLog
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
