package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type ConnectorFilter struct {
	// Kinds is an optional list of connector kinds to filter by.
	Kinds []model.ConnectorKind `json:"kinds"`
}

type ConnectorInfo struct {
	ConnectorId   string                `json:"connector_id"`
	Description   string                `json:"description"`
	Kind          string                `json:"kind"`
	DeploymentIds []string              `json:"deployment_ids,omitempty"`
	IsActive      bool                  `json:"is_active"`
	Config        model.ConnectorConfig `json:"config"`
}

type ConnectorCursor func(info *ConnectorInfo, hasMore bool) error

func (c *client) ListConnectors(filter ConnectorFilter, cursor ConnectorCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("CONNECTOR.LIST"), filter, func(b []byte, hasMore bool) error {
		if b == nil {
			return cursor(nil, hasMore)
		}

		var info ConnectorInfo
		if err := json.Unmarshal(b, &info); err != nil {
			return err
		}
		return cursor(&info, hasMore)
	}, opts...)
}
