package builders

import "github.com/synadia-io/connect/model"

type ProducerStepCoreBuilder struct {
	res *model.ProducerStepCore
}

func ProducerStepCore(subject string) *ProducerStepCoreBuilder {
	return &ProducerStepCoreBuilder{
		res: &model.ProducerStepCore{
			Subject: subject,
		},
	}
}

func (b *ProducerStepCoreBuilder) Build() model.ProducerStepCore {
	return *b.res
}
