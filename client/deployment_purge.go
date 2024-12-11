package client

import "encoding/json"

type DeploymentPurgeInfo struct {
	Error string `json:"error,omitempty"`

	ConnectorId  string `json:"connector_id"`
	DeploymentId string `json:"deployment_id"`
}

type DeploymentPurgeCursor func(item *DeploymentPurgeInfo, hasMore bool) error

func (c *client) PurgeDeployments(filter DeploymentFilter, cursor DeploymentPurgeCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("DEPLOYMENT.PURGE"), filter, func(b []byte, hasMore bool) error {
		if b == nil {
			return cursor(nil, hasMore)
		} else {
			var resp DeploymentPurgeInfo
			if err := json.Unmarshal(b, &resp); err != nil {
				return err
			}
			return cursor(&resp, hasMore)
		}
	}, opts...)
}
