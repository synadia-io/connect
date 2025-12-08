package builders

import (
	"github.com/synadia-io/connect/spec"
)

type ConsumerStepCoreBuilder struct {
	res *spec.ConsumerStepSpecCore
}

func ConsumerStepCore(subject string) *ConsumerStepCoreBuilder {
	return &ConsumerStepCoreBuilder{
		res: &spec.ConsumerStepSpecCore{
			Subject: subject,
		},
	}
}

func (b *ConsumerStepCoreBuilder) Queue(v string) *ConsumerStepCoreBuilder {
	b.res.Queue = &v
	return b
}

func (b *ConsumerStepCoreBuilder) Build() spec.ConsumerStepSpecCore {
	return *b.res
}
