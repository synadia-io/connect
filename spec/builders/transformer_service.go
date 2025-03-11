package builders

import (
    "github.com/synadia-io/connect/spec"
)

type ServiceTransformerStepBuilder struct {
    res *spec.TransformerStepSpecService
}

func ServiceTransformerStep(endpoint string, nats *NatsConfigBuilder) *ServiceTransformerStepBuilder {
    return &ServiceTransformerStepBuilder{
        res: &spec.TransformerStepSpecService{
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

func (b *ServiceTransformerStepBuilder) Build() spec.TransformerStepSpecService {
    return *b.res
}
