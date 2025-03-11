package builders

import (
    "github.com/synadia-io/connect/spec"
)

type MappingTransformerStepBuilder struct {
    res *spec.TransformerStepSpecMapping
}

func MappingTransformerStep(sourceCode string) *MappingTransformerStepBuilder {
    return &MappingTransformerStepBuilder{
        res: &spec.TransformerStepSpecMapping{
            Sourcecode: sourceCode,
        },
    }
}

func (b *MappingTransformerStepBuilder) Build() spec.TransformerStepSpecMapping {
    return *b.res
}
