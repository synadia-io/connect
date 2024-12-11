package library

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	getRuntimeRequest struct {
		RuntimeId string `json:"runtime_id"`
	}

	getRuntimeResponse struct {
		Found    bool           `json:"found"`
		Revision uint64         `json:"revision"`
		Runtime  *model.Runtime `json:"runtime,omitempty"`
	}
)

func (c *client) GetRuntime(runtimeId string, opts ...Opt) (*model.Runtime, error) {
	req := getRuntimeRequest{RuntimeId: runtimeId}

	b, err := c.Request(c.serviceSubject("RUNTIME.GET"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp getRuntimeResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return resp.Runtime, nil
}
