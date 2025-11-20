package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type TransformerStepBuilder struct {
	res *spec.TransformerStepSpec
}

func TransformerStep() *TransformerStepBuilder {
	return &TransformerStepBuilder{
		res: &spec.TransformerStepSpec{},
	}
}

func (b *TransformerStepBuilder) Composite(v *CompositeTransformerStepBuilder) *TransformerStepBuilder {
	r := v.Build()
	b.res.Composite = &r
	return b
}

func (b *TransformerStepBuilder) Mapping(v *MappingTransformerStepBuilder) *TransformerStepBuilder {
	r := v.Build()
	b.res.Mapping = &r
	return b
}

func (b *TransformerStepBuilder) Service(v *ServiceTransformerStepBuilder) *TransformerStepBuilder {
	r := v.Build()
	b.res.Service = &r
	return b
}

func (b *TransformerStepBuilder) Explode(v *ExplodeTransformerStepBuilder) *TransformerStepBuilder {
	r := v.Build()
	b.res.Explode = &r
	return b
}

func (b *TransformerStepBuilder) Combine(v *CombineTransformerStepBuilder) *TransformerStepBuilder {
	r := v.Build()
	b.res.Combine = &r
	return b
}

func (b *TransformerStepBuilder) Build() spec.TransformerStepSpec {
	return *b.res
}
