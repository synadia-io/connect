package builders

import "github.com/synadia-io/connect/model"

type ProducerStepBuilder struct {
    res *model.ProducerStep
}

func ProducerStep(nats *NatsConfigBuilder) *ProducerStepBuilder {
    return &ProducerStepBuilder{
        res: &model.ProducerStep{
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

func (b *ProducerStepBuilder) Build() model.ProducerStep {
    return *b.res
}
