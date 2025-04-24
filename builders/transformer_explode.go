package builders

import "github.com/synadia-io/connect/model"

type ExplodeTransformerStepBuilder struct {
    res *model.ExplodeTransformerStep
}

func ExplodeTransformerStep() *ExplodeTransformerStepBuilder {
    return &ExplodeTransformerStepBuilder{}
}

func (b *ExplodeTransformerStepBuilder) Build() model.ExplodeTransformerStep {
    return *b.res
}
