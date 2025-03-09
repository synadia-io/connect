package builders

import "github.com/synadia-io/connect/model"

type ConsumerStepBuilder struct {
    res *model.ConsumerStep
}

func ConsumerStep(nats *NatsConfigBuilder) *ConsumerStepBuilder {
    return &ConsumerStepBuilder{
        res: &model.ConsumerStep{
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

func (b *ConsumerStepBuilder) Build() model.ConsumerStep {
    return *b.res
}
