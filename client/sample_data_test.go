package client_test

import (
	im "github.com/synadia-io/connect-node/model"
	"github.com/synadia-io/connect/model"
)

func internalInletConfig() im.ConnectorConfig {
	return im.ConnectorConfig{
		Description: "A test inlet",
		Workload:    "synadia.com/vanilla-runtime:latest",
		Metrics: &im.MetricsEndpoint{
			Port: 4123,
			Path: "/metrics",
		},
		Steps: &im.Steps{
			Source:   &im.Source{},
			Producer: &im.Producer{},
		},
	}
}

func publicInletConfig() model.ConnectorConfig {
	return model.ConnectorConfig{
		Description: "A test inlet",
		Workload:    "synadia.com/vanilla-runtime:latest",
		Metrics: &model.MetricsEndpoint{
			Port: 4123,
			Path: "/metrics",
		},
		Steps: &model.Steps{
			Source:   &model.Source{},
			Producer: &model.Producer{},
		},
	}
}
