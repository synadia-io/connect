package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	updateConnectorRequest struct {
		ConnectorId string                `json:"connector_id"`
		Updates     model.ConnectorConfig `json:"updates"`
	}
	updateConnectorResponse struct {
		model.ConnectorConfig

		ConnectorId   string   `json:"connector_id"`
		Revision      uint64   `json:"revision"`
		DeploymentIds []string `json:"deployment_ids,omitempty"`
	}
)

func (c *client) UpdateConnector(id string, updates model.ConnectorConfig, opts ...Opt) (*model.Connector, error) {
	req := updateConnectorRequest{
		ConnectorId: id,
		Updates:     updates,
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.UPDATE"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp updateConnectorResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &model.Connector{
		ConnectorConfig: resp.ConnectorConfig,
		Id:              resp.ConnectorId,
		DeploymentIds:   resp.DeploymentIds,
	}, nil
}
