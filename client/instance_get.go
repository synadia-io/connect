package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	getInstanceRequest struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id"`
		InstanceId   string `json:"instance_id"`
	}
	getInstanceResponse struct {
		Found    bool   `json:"found"`
		Revision uint64 `json:"revision"`

		model.Instance
	}
)

func (c *client) GetInstance(connectorId string, deploymentId string, instanceId string, opts ...Opt) (*model.Instance, error) {
	req := getInstanceRequest{
		ConnectorId:  connectorId,
		DeploymentId: deploymentId,
		InstanceId:   instanceId,
	}

	b, err := c.Request(c.serviceSubject("INSTANCE.GET"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp getInstanceResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return &resp.Instance, nil
}
