package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	getConnectorRequest struct {
		Id string `json:"id"`
	}
	getConnectorResponse struct {
		model.ConnectorConfig

		Found    bool   `json:"found"`
		Revision uint64 `json:"revision"`

		Id            string   `json:"id"`
		DeploymentIds []string `json:"deployment_ids,omitempty"`
	}
)

func (c *client) GetConnector(id string, opts ...Opt) (*model.Connector, error) {
	req := getConnectorRequest{Id: id}

	b, err := c.Request(c.serviceSubject("CONNECTOR.GET"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp getConnectorResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	var result *model.Connector
	if resp.Found {
		result = &model.Connector{
			ConnectorConfig: resp.ConnectorConfig,
			Id:              resp.Id,
			DeploymentIds:   resp.DeploymentIds,
		}
	}

	return result, nil
}
