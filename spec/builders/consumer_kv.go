package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type ConsumerStepKvBuilder struct {
	res *spec.ConsumerStepSpecKv
}

func ConsumerStepKv(bucket string, key string) *ConsumerStepKvBuilder {
	return &ConsumerStepKvBuilder{
		res: &spec.ConsumerStepSpecKv{
			Bucket: bucket,
			Key:    key,
		},
	}
}

func (b *ConsumerStepKvBuilder) Build() spec.ConsumerStepSpecKv {
	return *b.res
}
