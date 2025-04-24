package builders

import (
    "github.com/synadia-io/connect/spec"
)

type ExplodeTransformerStepBuilder struct {
    res *spec.TransformerStepSpecExplode
}

func ExplodeTransformerStep() *ExplodeTransformerStepBuilder {
    return &ExplodeTransformerStepBuilder{}
}

func (b *ExplodeTransformerStepBuilder) Build() spec.TransformerStepSpecExplode {
    return *b.res
}
