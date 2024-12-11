package library

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	getComponentRequest struct {
		RuntimeId string              `json:"runtime_id"`
		VersionId string              `json:"version_id"`
		Kind      model.ComponentKind `json:"kind"`
		Name      string              `json:"name"`
	}

	getComponentResponse struct {
		Found     bool             `json:"found"`
		Revision  uint64           `json:"revision"`
		Component *model.Component `json:"component,omitempty"`
	}
)

func (c *client) GetComponent(runtimeId string, versionId string, kind model.ComponentKind, name string, opts ...Opt) (*model.Component, error) {
	req := getComponentRequest{
		RuntimeId: runtimeId,
		VersionId: versionId,
		Kind:      kind,
		Name:      name,
	}

	b, err := c.Request(c.serviceSubject("COMPONENT.GET"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp getComponentResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return resp.Component, nil
}
