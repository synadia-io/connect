package builders

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StepsBuilder", func() {
	var builder *StepsBuilder

	BeforeEach(func() {
		builder = Steps()
	})

	Describe("Steps", func() {
		It("should create a new steps builder", func() {
			Expect(builder).ToNot(BeNil())
			Expect(builder.steps).ToNot(BeNil())
		})
	})

	Describe("Producer", func() {
		It("should add a producer step", func() {
			natsConfig := NatsConfig().Url("nats://localhost:4222")
			producer := ProducerStep(natsConfig)
			builder.Producer(producer)

			steps := builder.Build()
			Expect(steps.Producer).ToNot(BeNil())
		})
	})

	Describe("Consumer", func() {
		It("should add a consumer step", func() {
			natsConfig := NatsConfig().Url("nats://localhost:4222")
			consumer := ConsumerStep(natsConfig)
			builder.Consumer(consumer)

			steps := builder.Build()
			Expect(steps.Consumer).ToNot(BeNil())
		})
	})

	Describe("Source", func() {
		It("should add a source step", func() {
			source := SourceStep("test-source")
			builder.Source(source)

			steps := builder.Build()
			Expect(steps.Source).ToNot(BeNil())
			Expect(steps.Source.Type).To(Equal("test-source"))
		})
	})

	Describe("Sink", func() {
		It("should add a sink step", func() {
			sink := SinkStep("test-sink")
			builder.Sink(sink)

			steps := builder.Build()
			Expect(steps.Sink).ToNot(BeNil())
			Expect(steps.Sink.Type).To(Equal("test-sink"))
		})
	})

	Describe("Transformer", func() {
		It("should add a transformer step", func() {
			transformer := TransformerStep()
			builder.Transformer(transformer)

			steps := builder.Build()
			Expect(steps.Transformer).ToNot(BeNil())
		})
	})

	Describe("Build", func() {
		It("should build a complete steps configuration", func() {
			steps := builder.
				Source(SourceStep("generate")).
				Sink(SinkStep("http")).
				Transformer(TransformerStep().Mapping(MappingTransformerStep("root = this"))).
				Build()

			Expect(steps.Source).ToNot(BeNil())
			Expect(steps.Sink).ToNot(BeNil())
			Expect(steps.Transformer).ToNot(BeNil())
		})

		It("should build steps with only source", func() {
			steps := builder.
				Source(SourceStep("generate")).
				Build()

			Expect(steps.Source).ToNot(BeNil())
			Expect(steps.Sink).To(BeNil())
			Expect(steps.Producer).To(BeNil())
			Expect(steps.Consumer).To(BeNil())
		})
	})

	Describe("Fluent API", func() {
		It("should support method chaining", func() {
			result := builder.
				Source(SourceStep("test")).
				Sink(SinkStep("output")).
				Transformer(TransformerStep())

			Expect(result).To(Equal(builder))
		})
	})
})
