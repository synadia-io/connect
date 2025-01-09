package client

import (
	"encoding/json"

	"github.com/synadia-io/connect/model"
)

type (
	getConnectorEventsRequest struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id,omitempty"`
		InstanceId   string `json:"instance_id,omitempty"`
	}
)

func (c *client) GetEvents(connectorId, deploymentId, instanceId string, opts ...Opt) ([]*model.InstanceEvent, error) {
	req := getConnectorEventsRequest{
		ConnectorId:  connectorId,
		DeploymentId: deploymentId,
		InstanceId:   instanceId,
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.EVENTS"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp []*model.InstanceEvent
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
