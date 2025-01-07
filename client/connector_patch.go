package client

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/synadia-io/connect/model"
)

type (
	patchConnectorRequest struct {
		ConnectorId string `json:"connector_id"`
		Patch       []byte `json:"patch"`
	}
	patchConnectorResponse struct {
		model.ConnectorConfig

		ConnectorId   string   `json:"connector_id"`
		Revision      uint64   `json:"revision"`
		DeploymentIds []string `json:"deployment_ids,omitempty"`
	}
)

func (c *client) PatchConnector(id string, updates model.ConnectorConfig, opts ...Opt) (*model.Connector, error) {
	connector, err := c.GetConnector(id, opts...)
	if err != nil {
		return nil, err
	}

	orig, err := json.Marshal(connector.ConnectorConfig)
	if err != nil {
		return nil, err
	}

	updated, err := json.Marshal(updates)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreateMergePatch(orig, updated)
	if err != nil {
		return nil, err
	}

	req := patchConnectorRequest{
		ConnectorId: id,
		Patch:       patch,
	}
	b, err := c.Request(c.serviceSubject("CONNECTOR.PATCH"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp patchConnectorResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &model.Connector{
		ConnectorConfig: resp.ConnectorConfig,
		Id:              resp.ConnectorId,
		DeploymentIds:   resp.DeploymentIds,
	}, nil
}
