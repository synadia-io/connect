package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type SinkStepBuilder struct {
	res *spec.SinkStepSpec
}

func SinkStep(kind string) *SinkStepBuilder {
	return &SinkStepBuilder{
		res: &spec.SinkStepSpec{
			Type: kind,
		},
	}
}

func (b *SinkStepBuilder) SetString(key string, value string) *SinkStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SinkStepBuilder) SetStrings(key string, value ...string) *SinkStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SinkStepBuilder) SetInt(key string, value int) *SinkStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SinkStepBuilder) SetBool(key string, value bool) *SinkStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SinkStepBuilder) Build() spec.SinkStepSpec {
	return *b.res
}
