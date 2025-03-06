package builders

import "github.com/synadia-io/connect/v2/model"

type TransformerStepBuilder struct {
    res *model.TransformerStep
}

func TransformerStep() *TransformerStepBuilder {
    return &TransformerStepBuilder{
        res: &model.TransformerStep{},
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

func (b *TransformerStepBuilder) Build() model.TransformerStep {
    return *b.res
}
