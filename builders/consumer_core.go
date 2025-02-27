package builders

import "github.com/synadia-io/connect/model"

type ConsumerStepCoreBuilder struct {
    res *model.ConsumerStepCore
}

func ConsumerStepCore(subject string) *ConsumerStepCoreBuilder {
    return &ConsumerStepCoreBuilder{
        res: &model.ConsumerStepCore{
            Subject: subject,
        },
    }
}

func (b *ConsumerStepCoreBuilder) Queue(v string) *ConsumerStepCoreBuilder {
    b.res.Queue = &v
    return b
}

func (b *ConsumerStepCoreBuilder) Build() model.ConsumerStepCore {
    return *b.res
}
