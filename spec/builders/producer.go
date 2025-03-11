package builders

import (
    "github.com/synadia-io/connect/spec"
)

type ProducerStepBuilder struct {
    res *spec.ProducerStepSpec
}

func ProducerStep(nats *NatsConfigBuilder) *ProducerStepBuilder {
    return &ProducerStepBuilder{
        res: &spec.ProducerStepSpec{
            Nats:    nats.Build(),
            Threads: 1,
        },
    }
}

func (b *ProducerStepBuilder) Core(v *ProducerStepCoreBuilder) *ProducerStepBuilder {
    r := v.Build()
    b.res.Core = &r
    return b
}

func (b *ProducerStepBuilder) Stream(v *ProducerStepStreamBuilder) *ProducerStepBuilder {
    r := v.Build()
    b.res.Stream = &r
    return b
}

func (b *ProducerStepBuilder) Kv(v *ProducerStepKvBuilder) *ProducerStepBuilder {
    r := v.Build()
    b.res.Kv = &r
    return b
}

func (b *ProducerStepBuilder) Build() spec.ProducerStepSpec {
    return *b.res
}
