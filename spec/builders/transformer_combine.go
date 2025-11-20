package builders

import "github.com/synadia-io/connect/v2/spec"

type CombineTransformerStepBuilder struct {
	res *spec.TransformerStepSpecCombine
}

func CombineTransformerStep() *CombineTransformerStepBuilder {
	return &CombineTransformerStepBuilder{
		res: &spec.TransformerStepSpecCombine{
			Format: spec.TransformerStepSpecCombineFormatJsonArray,
		},
	}
}

func (b *CombineTransformerStepBuilder) Format(format spec.TransformerStepSpecCombineFormat) *CombineTransformerStepBuilder {
	b.res.Format = format
	return b
}

func (b *CombineTransformerStepBuilder) Path(path string) *CombineTransformerStepBuilder {
	b.res.Path = path
	return b
}

func (b *CombineTransformerStepBuilder) Build() spec.TransformerStepSpecCombine {
	return *b.res
}
