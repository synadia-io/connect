package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type InstanceFilter struct {
	ConnectorId  string `json:"connector_id,omitempty"`
	DeploymentId string `json:"deployment_id,omitempty"`
}

type InstanceInfo struct {
	ConnectorId  string               `json:"connector_id"`
	DeploymentId string               `json:"deployment_id"`
	InstanceId   string               `json:"instance_id"`
	Status       model.InstanceStatus `json:"status"`
	Uptime       string               `json:"uptime,omitempty"`
}

type InstanceCursor func(item *InstanceInfo, hasMore bool) error

func (c *client) ListInstances(filter InstanceFilter, cursor InstanceCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("INSTANCE.LIST"), filter, func(b []byte, hasMore bool) error {
		var info *InstanceInfo

		if b != nil {
			var resp InstanceInfo
			if err := json.Unmarshal(b, &resp); err != nil {
				return err
			}
		}

		return cursor(info, hasMore)
	}, opts...)
}
