package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	createConnectorRequest struct {
		ConnectorId string                `json:"connector_id"`
		Config      model.ConnectorConfig `json:"config"`
	}
	createConnectorResponse struct {
		model.ConnectorConfig

		ConnectorId string `json:"connector_id"`
		Revision    uint64 `json:"revision"`
	}
)

func (c *client) CreateConnector(id string, config model.ConnectorConfig, opts ...Opt) (*model.Connector, error) {
	req := createConnectorRequest{
		ConnectorId: id,
		Config:      config,
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.CREATE"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp createConnectorResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &model.Connector{
		ConnectorConfig: resp.ConnectorConfig,
		Id:              resp.ConnectorId,
	}, nil
}
