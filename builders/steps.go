package builders

import "github.com/synadia-io/connect/v2/model"

type StepsBuilder struct {
	steps *model.Steps
}

func Steps() *StepsBuilder {
	return &StepsBuilder{
		steps: &model.Steps{},
	}
}

func (b *StepsBuilder) Producer(v *ProducerStepBuilder) *StepsBuilder {
	r := v.Build()
	b.steps.Producer = &r
	return b
}

func (b *StepsBuilder) Consumer(v *ConsumerStepBuilder) *StepsBuilder {
	r := v.Build()
	b.steps.Consumer = &r
	return b
}

func (b *StepsBuilder) Source(v *SourceStepBuilder) *StepsBuilder {
	r := v.Build()
	b.steps.Source = &r
	return b
}

func (b *StepsBuilder) Sink(v *SinkStepBuilder) *StepsBuilder {
	r := v.Build()
	b.steps.Sink = &r
	return b
}

func (b *StepsBuilder) Transformer(v *TransformerStepBuilder) *StepsBuilder {
	r := v.Build()
	b.steps.Transformer = &r
	return b
}

func (b *StepsBuilder) Build() model.Steps {
	return *b.steps
}
