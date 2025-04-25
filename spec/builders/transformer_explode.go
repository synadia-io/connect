package builders

import (
    "github.com/synadia-io/connect/spec"
)

type ExplodeTransformerStepBuilder struct {
    res *spec.TransformerStepSpecExplode
}

func ExplodeTransformerStep() *ExplodeTransformerStepBuilder {
    return &ExplodeTransformerStepBuilder{
        res: &spec.TransformerStepSpecExplode{
            Format: spec.TransformerStepSpecExplodeFormatJsonArray,
        },
    }
}

func (b *ExplodeTransformerStepBuilder) Format(format spec.TransformerStepSpecExplodeFormat) *ExplodeTransformerStepBuilder {
    b.res.Format = format
    return b
}

func (b *ExplodeTransformerStepBuilder) Delimiter(delimiter string) *ExplodeTransformerStepBuilder {
    b.res.Delimiter = delimiter
    return b
}

func (b *ExplodeTransformerStepBuilder) Build() spec.TransformerStepSpecExplode {
    return *b.res
}
