package client

import "encoding/json"

type (
	undeployConnectorRequest struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id"`
		Timeout      string `json:"timeout"`
		Force        bool   `json:"force"`
	}

	ConnectorUndeployResult struct {
		DeploymentId   string            `json:"deployment_id"`
		InstanceErrors map[string]string `json:"errors,omitempty"`
		Error          string            `json:"error,omitempty"`
	}
)

type UndeployOpt func(*undeployConnectorRequest)

func WithUndeployTimeout(timeout string) UndeployOpt {
	return func(r *undeployConnectorRequest) {
		r.Timeout = timeout
	}
}

func WithUndeployForce(force bool) UndeployOpt {
	return func(r *undeployConnectorRequest) {
		r.Force = force
	}
}

func WithDeploymentId(deploymentId string) UndeployOpt {
	return func(r *undeployConnectorRequest) {
		r.DeploymentId = deploymentId
	}
}

func (c *client) UndeployConnector(connectorId string, opts ...UndeployOpt) (*ConnectorUndeployResult, error) {
	req := undeployConnectorRequest{
		ConnectorId: connectorId,
	}

	for _, opt := range opts {
		opt(&req)
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.UNDEPLOY"), req)
	if err != nil {
		return nil, err
	}

	var resp ConnectorUndeployResult
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
