package model

type (
	Deployment struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id"`

		DeploymentConfig

		Status      DeploymentStatus `json:"status"`
		Allocations map[string]int   `json:"allocations,omitempty"`
	}

	DeploymentConfig struct {
		PlacementTags []string `json:"placement_tags,omitempty"`
		Replicas      int      `json:"replicas"`

		Workload DeploymentWorkload `json:"workload"`

		Targets map[string]DeployTarget `json:"targets"`
	}

	DeploymentWorkload struct {
		Location string           `json:"location"`
		Config   []byte           `json:"config"`
		Metrics  *MetricsEndpoint `json:"metrics,omitempty"`
	}

	DeployTarget struct {
		AgentId    string `json:"agent_id"`
		InstanceId string `json:"instance_id"`
	}

	DeploymentStatus struct {
		Pending int `json:"pending"`
		Running int `json:"running"`
		Errored int `json:"errored"`
		Stopped int `json:"stopped"`
	}
)
