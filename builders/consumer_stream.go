package builders

import "github.com/synadia-io/connect/model"

type ConsumerStepStreamBuilder struct {
	res *model.ConsumerStepStream
}

func ConsumerStepStream(subject string) *ConsumerStepStreamBuilder {
	return &ConsumerStepStreamBuilder{
		res: &model.ConsumerStepStream{
			Subject: subject,
		},
	}
}

func (b *ConsumerStepStreamBuilder) Build() model.ConsumerStepStream {
	return *b.res
}
