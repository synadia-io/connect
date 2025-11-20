package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type SourceStepBuilder struct {
	res *spec.SourceStepSpec
}

func SourceStep(kind string) *SourceStepBuilder {
	return &SourceStepBuilder{
		res: &spec.SourceStepSpec{
			Type: kind,
		},
	}
}

func (b *SourceStepBuilder) SetString(key string, value string) *SourceStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SourceStepBuilder) SetStrings(key string, value ...string) *SourceStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SourceStepBuilder) SetInt(key string, value int) *SourceStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SourceStepBuilder) SetBool(key string, value bool) *SourceStepBuilder {
	if b.res.Config == nil {
		b.res.Config = make(map[string]interface{})
	}
	b.res.Config[key] = value
	return b
}

func (b *SourceStepBuilder) Build() spec.SourceStepSpec {
	return *b.res
}
