package client

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	deployConnectorRequest struct {
		ConnectorId   string            `json:"connector_id"`
		Timeout       string            `json:"timeout"`
		EnvVars       map[string]string `json:"env_vars"`
		PlacementTags []string          `json:"placement_tags"`
		Replicas      int               `json:"replicas"`
		Pull          bool              `json:"pull"`
	}

	ConnectorDeployResult struct {
		InstanceErrors map[string]string `json:"errors,omitempty"`
		Error          string            `json:"error,omitempty"`

		DeploymentId string                        `json:"deployment_id"`
		Targets      map[string]model.DeployTarget `json:"targets"`
	}
)

type DeployOpt func(*deployConnectorRequest)

func WithDeployTimeout(timeout string) DeployOpt {
	return func(r *deployConnectorRequest) {
		r.Timeout = timeout
	}
}

func WithDeployEnvVars(envVars map[string]string) DeployOpt {
	return func(r *deployConnectorRequest) {
		if r.EnvVars == nil {
			r.EnvVars = make(map[string]string)
		}

		for k, v := range envVars {
			r.EnvVars[k] = v
		}
	}
}

func WithDeployEnvVar(key string, value string) DeployOpt {
	return func(request *deployConnectorRequest) {
		if request.EnvVars == nil {
			request.EnvVars = make(map[string]string)
		}
		request.EnvVars[key] = value
	}
}

func WithDeployPlacementTags(tags []string) DeployOpt {
	return func(r *deployConnectorRequest) {
		r.PlacementTags = tags
	}
}

func WithDeployReplicas(replicas int) DeployOpt {
	return func(r *deployConnectorRequest) {
		r.Replicas = replicas
	}
}

func WithPull(pull bool) DeployOpt {
	return func(r *deployConnectorRequest) {
		r.Pull = pull
	}
}

func (c *client) DeployConnector(id string, opts ...DeployOpt) (*ConnectorDeployResult, error) {
	req := deployConnectorRequest{ConnectorId: id}

	for _, opt := range opts {
		opt(&req)
	}

	b, err := c.Request(c.serviceSubject("CONNECTOR.DEPLOY"), req)
	if err != nil {
		return nil, err
	}

	var resp ConnectorDeployResult
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
