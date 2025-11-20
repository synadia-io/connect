package builders

import (
	"github.com/synadia-io/connect/v2/spec"
)

type CompositeTransformerStepBuilder struct {
	res *spec.TransformerStepSpecComposite
}

func CompositeTransformerStep() *CompositeTransformerStepBuilder {
	return &CompositeTransformerStepBuilder{
		res: &spec.TransformerStepSpecComposite{},
	}
}

func (b *CompositeTransformerStepBuilder) Sequential(v ...*TransformerStepBuilder) *CompositeTransformerStepBuilder {
	for _, t := range v {
		b.res.Sequential = append(b.res.Sequential, t.Build())
	}
	return b
}

func (b *CompositeTransformerStepBuilder) Build() spec.TransformerStepSpecComposite {
	return *b.res
}
