package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type ProducerStepStreamBuilder struct {
	res *spec.ProducerStepSpecStream
}

func ProducerStepStream(subject string) *ProducerStepStreamBuilder {
	return &ProducerStepStreamBuilder{
		res: &spec.ProducerStepSpecStream{
			Subject: subject,
		},
	}
}

func (b *ProducerStepStreamBuilder) Build() spec.ProducerStepSpecStream {
	return *b.res
}
