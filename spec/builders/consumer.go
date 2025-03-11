package builders

import (
    "github.com/synadia-io/connect/spec"
)

type ConsumerStepBuilder struct {
    res *spec.ConsumerStepSpec
}

func ConsumerStep(nats *NatsConfigBuilder) *ConsumerStepBuilder {
    return &ConsumerStepBuilder{
        res: &spec.ConsumerStepSpec{
            Nats: nats.Build(),
        },
    }
}

func (b *ConsumerStepBuilder) Core(v *ConsumerStepCoreBuilder) *ConsumerStepBuilder {
    c := v.Build()
    b.res.Core = &c
    return b
}

func (b *ConsumerStepBuilder) Kv(v *ConsumerStepKvBuilder) *ConsumerStepBuilder {
    c := v.Build()
    b.res.Kv = &c
    return b
}

func (b *ConsumerStepBuilder) Stream(v *ConsumerStepStreamBuilder) *ConsumerStepBuilder {
    c := v.Build()
    b.res.Stream = &c
    return b
}

func (b *ConsumerStepBuilder) Build() spec.ConsumerStepSpec {
    return *b.res
}
