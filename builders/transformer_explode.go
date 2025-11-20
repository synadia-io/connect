package builders

import "github.com/synadia-io/connect/v2/model"

type ExplodeTransformerStepBuilder struct {
	res *model.ExplodeTransformerStep
}

func ExplodeTransformerStep() *ExplodeTransformerStepBuilder {
	return &ExplodeTransformerStepBuilder{
		res: &model.ExplodeTransformerStep{
			Format: model.ExplodeTransformerStepFormatJsonArray,
		},
	}
}

func (b *ExplodeTransformerStepBuilder) Format(format model.ExplodeTransformerStepFormat) *ExplodeTransformerStepBuilder {
	b.res.Format = format
	return b
}

func (b *ExplodeTransformerStepBuilder) Delimiter(delimiter string) *ExplodeTransformerStepBuilder {
	b.res.Delimiter = delimiter
	return b
}

func (b *ExplodeTransformerStepBuilder) Build() model.ExplodeTransformerStep {
	return *b.res
}
