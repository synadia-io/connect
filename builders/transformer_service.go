package builders

import "github.com/synadia-io/connect/v2/model"

type ServiceTransformerStepBuilder struct {
	res *model.ServiceTransformerStep
}

func ServiceTransformerStep(endpoint string, nats *NatsConfigBuilder) *ServiceTransformerStepBuilder {
	return &ServiceTransformerStepBuilder{
		res: &model.ServiceTransformerStep{
			Endpoint: endpoint,
			Nats:     nats.Build(),
			Timeout:  "5s",
		},
	}
}

func (b *ServiceTransformerStepBuilder) Timeout(timeout string) *ServiceTransformerStepBuilder {
	b.res.Timeout = timeout
	return b
}

func (b *ServiceTransformerStepBuilder) Build() model.ServiceTransformerStep {
	return *b.res
}
