package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type DeploymentFilter struct {
	ConnectorId    string                 `json:"connector_id"`
	InstanceStatus []model.InstanceStatus `json:"instance_status"`
}

type DeploymentInfo struct {
	ConnectorId    string                  `json:"connector_id"`
	DeploymentId   string                  `json:"deployment_id"`
	Replicas       int                     `json:"replicas"`
	PlacementTags  []string                `json:"placement_tags"`
	Workload       string                  `json:"workload"`
	MetricsEnabled bool                    `json:"metrics_enabled"`
	Status         *model.DeploymentStatus `json:"status"`
}

type DeploymentCursor func(item *DeploymentInfo, hasMore bool) error

func (c *client) ListDeployments(filter DeploymentFilter, cursor DeploymentCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("DEPLOYMENT.LIST"), filter, func(b []byte, hasMore bool) error {
		if b == nil {
			return cursor(nil, hasMore)
		} else {
			var resp DeploymentInfo
			if err := json.Unmarshal(b, &resp); err != nil {
				return err
			}
			return cursor(&resp, hasMore)
		}
	}, opts...)
}
