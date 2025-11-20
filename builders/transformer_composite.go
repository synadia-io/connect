package builders

import "github.com/synadia-io/connect/v2/model"

type CompositeTransformerStepBuilder struct {
	res *model.CompositeTransformerStep
}

func CompositeTransformerStep() *CompositeTransformerStepBuilder {
	return &CompositeTransformerStepBuilder{
		res: &model.CompositeTransformerStep{},
	}
}

func (b *CompositeTransformerStepBuilder) Sequential(v ...*TransformerStepBuilder) *CompositeTransformerStepBuilder {
	for _, t := range v {
		b.res.Sequential = append(b.res.Sequential, t.Build())
	}
	return b
}

func (b *CompositeTransformerStepBuilder) Build() model.CompositeTransformerStep {
	return *b.res
}
