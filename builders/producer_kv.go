package builders

import "github.com/synadia-io/connect/v2/model"

type ProducerStepKvBuilder struct {
	res *model.ProducerStepKv
}

func ProducerStepKv(bucket string, key string) *ProducerStepKvBuilder {
	return &ProducerStepKvBuilder{
		res: &model.ProducerStepKv{
			Bucket: bucket,
			Key:    key,
		},
	}
}

func (b *ProducerStepKvBuilder) Build() model.ProducerStepKv {
	return *b.res
}
