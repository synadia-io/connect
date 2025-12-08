package builders

import "github.com/synadia-io/connect/model"

type CombineTransformerStepBuilder struct {
	res *model.CombineTransformerStep
}

func CombineTransformerStep() *CombineTransformerStepBuilder {
	return &CombineTransformerStepBuilder{
		res: &model.CombineTransformerStep{
			Format: model.CombineTransformerStepFormatJsonArray,
		},
	}
}

func (b *CombineTransformerStepBuilder) Format(format model.CombineTransformerStepFormat) *CombineTransformerStepBuilder {
	b.res.Format = format
	return b
}

func (b *CombineTransformerStepBuilder) Path(path string) *CombineTransformerStepBuilder {
	b.res.Path = path
	return b
}

func (b *CombineTransformerStepBuilder) Build() model.CombineTransformerStep {
	return *b.res
}
