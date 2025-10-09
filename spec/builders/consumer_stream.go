package builders

import (
	"github.com/synadia-io/connect/spec"
)

type ConsumerStepStreamBuilder struct {
	res *spec.ConsumerStepSpecStream
}

func ConsumerStepStream(subject string) *ConsumerStepStreamBuilder {
	return &ConsumerStepStreamBuilder{
		res: &spec.ConsumerStepSpecStream{
			Subject: subject,
		},
	}
}

func (b *ConsumerStepStreamBuilder) Build() spec.ConsumerStepSpecStream {
	return *b.res
}
