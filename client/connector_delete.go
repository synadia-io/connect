package client

import "encoding/json"

type (
	deleteConnectorRequest struct {
		ConnectorId string `json:"connector_id"`
	}
	deleteConnectorResponse struct {
		Existed bool `json:"existed"`
	}
)

func (c *client) DeleteConnector(id string, opts ...Opt) (bool, error) {
	req := deleteConnectorRequest{ConnectorId: id}

	b, err := c.Request(c.serviceSubject("CONNECTOR.DELETE"), req, opts...)
	if err != nil {
		return false, err
	}

	var resp deleteConnectorResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return false, err
	}

	return resp.Existed, nil
}
