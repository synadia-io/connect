package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type ProducerStepCoreBuilder struct {
	res *spec.ProducerStepSpecCore
}

func ProducerStepCore(subject string) *ProducerStepCoreBuilder {
	return &ProducerStepCoreBuilder{
		res: &spec.ProducerStepSpecCore{
			Subject: subject,
		},
	}
}

func (b *ProducerStepCoreBuilder) Build() spec.ProducerStepSpecCore {
	return *b.res
}
