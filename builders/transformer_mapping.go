package builders

import "github.com/synadia-io/connect/v2/model"

type MappingTransformerStepBuilder struct {
    res *model.MappingTransformerStep
}

func MappingTransformerStep(sourceCode string) *MappingTransformerStepBuilder {
    return &MappingTransformerStepBuilder{
        res: &model.MappingTransformerStep{
            Sourcecode: sourceCode,
        },
    }
}

func (b *MappingTransformerStepBuilder) Build() model.MappingTransformerStep {
    return *b.res
}
