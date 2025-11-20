package builders

import "github.com/synadia-io/connect/v2/model"

type ProducerStepStreamBuilder struct {
	res *model.ProducerStepStream
}

func ProducerStepStream(subject string) *ProducerStepStreamBuilder {
	return &ProducerStepStreamBuilder{
		res: &model.ProducerStepStream{
			Subject: subject,
		},
	}
}

func (b *ProducerStepStreamBuilder) Build() model.ProducerStepStream {
	return *b.res
}
