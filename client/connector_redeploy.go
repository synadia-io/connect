package client

import "encoding/json"

type (
	redeployConnectorRequest struct {
		ConnectorId   string            `json:"connector_id"`
		Timeout       string            `json:"timeout"`
		EnvVars       map[string]string `json:"env_vars"`
		PlacementTags []string          `json:"placement_tags"`
		Replicas      int               `json:"replicas"`
		ForceUndeploy bool              `json:"force_undeploy"`
	}
	ConnectorRedeployResult struct {
		DeploymentId string                  `json:"deployment_id"`
		Deploy       ConnectorDeployResult   `json:"deploy"`
		Undeploy     ConnectorUndeployResult `json:"undeploy"`
	}
)

type RedeployOpt func(*redeployConnectorRequest)

func WithRedeployTimeout(timeout string) RedeployOpt {
	return func(r *redeployConnectorRequest) {
		r.Timeout = timeout
	}
}

func WithRedeployEnvVars(envVars map[string]string) RedeployOpt {
	return func(r *redeployConnectorRequest) {
		if r.EnvVars == nil {
			r.EnvVars = make(map[string]string)
		}

		for k, v := range envVars {
			r.EnvVars[k] = v
		}
	}
}

func WithRedeployEnvVar(key string, value string) RedeployOpt {
	return func(request *redeployConnectorRequest) {
		if request.EnvVars == nil {
			request.EnvVars = make(map[string]string)
		}
		request.EnvVars[key] = value
	}
}

func WithRedeployPlacementTags(tags []string) RedeployOpt {
	return func(r *redeployConnectorRequest) {
		r.PlacementTags = tags
	}
}

func WithRedeployReplicas(replicas int) RedeployOpt {
	return func(r *redeployConnectorRequest) {
		r.Replicas = replicas
	}
}

func WithRedeployForceUndeploy(forceUndeploy bool) RedeployOpt {
	return func(r *redeployConnectorRequest) {
		r.ForceUndeploy = forceUndeploy
	}
}

func (c *client) RedeployConnector(id string, opts ...RedeployOpt) (*ConnectorRedeployResult, error) {
	req := redeployConnectorRequest{
		ConnectorId: id,
	}

	for _, opt := range opts {
		opt(&req)
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.REDEPLOY"), req)
	if err != nil {
		return nil, err
	}

	var resp ConnectorRedeployResult
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
