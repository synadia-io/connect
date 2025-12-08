package builders

import "github.com/synadia-io/connect/model"

type SourceStepBuilder struct {
	res *model.SourceStep
}

func SourceStep(kind string) *SourceStepBuilder {
	return &SourceStepBuilder{
		res: &model.SourceStep{
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

func (b *SourceStepBuilder) Build() model.SourceStep {
	return *b.res
}
