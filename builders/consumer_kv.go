package builders

import "github.com/synadia-io/connect/v2/model"

type ConsumerStepKvBuilder struct {
	res *model.ConsumerStepKv
}

func ConsumerStepKv(bucket string, key string) *ConsumerStepKvBuilder {
	return &ConsumerStepKvBuilder{
		res: &model.ConsumerStepKv{
			Bucket: bucket,
			Key:    key,
		},
	}
}

func (b *ConsumerStepKvBuilder) Build() model.ConsumerStepKv {
	return *b.res
}
