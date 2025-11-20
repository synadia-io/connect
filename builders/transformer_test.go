package builders

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect/v2/model"
)

var _ = Describe("TransformerStepBuilder", func() {
	var builder *TransformerStepBuilder

	BeforeEach(func() {
		builder = TransformerStep()
	})

	Describe("TransformerStep", func() {
		It("should create a new transformer step builder", func() {
			Expect(builder).ToNot(BeNil())
			Expect(builder.res).ToNot(BeNil())
		})
	})

	Describe("Mapping", func() {
		It("should add a mapping transformer", func() {
			mapping := MappingTransformerStep("root.field = 'value'")
			builder.Mapping(mapping)

			result := builder.Build()
			Expect(result.Mapping).ToNot(BeNil())
			Expect(result.Mapping.Sourcecode).To(Equal("root.field = 'value'"))
		})

		It("should support method chaining", func() {
			result := builder.Mapping(MappingTransformerStep("test"))
			Expect(result).To(Equal(builder))
		})
	})

	Describe("Service", func() {
		It("should add a service transformer", func() {
			natsConfig := NatsConfig().Url("nats://localhost:4222")
			service := ServiceTransformerStep("service.transform", natsConfig)

			builder.Service(service)

			result := builder.Build()
			Expect(result.Service).ToNot(BeNil())
			Expect(result.Service.Endpoint).To(Equal("service.transform"))
			Expect(result.Service.Nats).ToNot(BeNil())
		})

		It("should support custom timeout", func() {
			natsConfig := NatsConfig()
			service := ServiceTransformerStep("service.endpoint", natsConfig).
				Timeout("10s")

			builder.Service(service)

			result := builder.Build()
			Expect(result.Service.Timeout).To(Equal("10s"))
		})
	})

	Describe("Composite", func() {
		It("should add a composite transformer", func() {
			sub1 := TransformerStep().Mapping(MappingTransformerStep("step1"))
			sub2 := TransformerStep().Mapping(MappingTransformerStep("step2"))

			composite := CompositeTransformerStep().Sequential(sub1, sub2)
			builder.Composite(composite)

			result := builder.Build()
			Expect(result.Composite).ToNot(BeNil())
			Expect(result.Composite.Sequential).To(HaveLen(2))
		})
	})

	Describe("Explode", func() {
		It("should add an explode transformer", func() {
			explode := ExplodeTransformerStep().
				Format(model.ExplodeTransformerStepFormatJsonArray)

			builder.Explode(explode)

			result := builder.Build()
			Expect(result.Explode).ToNot(BeNil())
			Expect(result.Explode.Format).To(Equal(model.ExplodeTransformerStepFormatJsonArray))
		})
	})

	Describe("Combine", func() {
		It("should add a combine transformer", func() {
			combine := CombineTransformerStep().
				Format(model.CombineTransformerStepFormatLines)

			builder.Combine(combine)

			result := builder.Build()
			Expect(result.Combine).ToNot(BeNil())
			Expect(result.Combine.Format).To(Equal(model.CombineTransformerStepFormatLines))
		})
	})

	Describe("Build", func() {
		It("should build an empty transformer", func() {
			result := builder.Build()

			Expect(result.Mapping).To(BeNil())
			Expect(result.Service).To(BeNil())
			Expect(result.Composite).To(BeNil())
			Expect(result.Explode).To(BeNil())
			Expect(result.Combine).To(BeNil())
		})

		It("should build with one transformer type", func() {
			result := builder.
				Mapping(MappingTransformerStep("root.name = 'test'")).
				Build()

			Expect(result.Mapping).ToNot(BeNil())
			Expect(result.Mapping.Sourcecode).To(Equal("root.name = 'test'"))
		})
	})

	Describe("Multiple transformer types", func() {
		It("should allow setting multiple transformer types", func() {
			// Note: In practice, only one transformer type should be used at a time
			// This test verifies the builder allows it technically
			result := builder.
				Mapping(MappingTransformerStep("mapping")).
				Service(ServiceTransformerStep("service", NatsConfig())).
				Explode(ExplodeTransformerStep()).
				Build()

			Expect(result.Mapping).ToNot(BeNil())
			Expect(result.Service).ToNot(BeNil())
			Expect(result.Explode).ToNot(BeNil())
		})
	})
})
