package builders

import (
    "github.com/synadia-io/connect/spec"
)

type ProducerStepKvBuilder struct {
    res *spec.ProducerStepSpecKv
}

func ProducerStepKv(bucket string, key string) *ProducerStepKvBuilder {
    return &ProducerStepKvBuilder{
        res: &spec.ProducerStepSpecKv{
            Bucket: bucket,
            Key:    key,
        },
    }
}

func (b *ProducerStepKvBuilder) Build() spec.ProducerStepSpecKv {
    return *b.res
}
