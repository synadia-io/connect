package builders

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SinkStepBuilder", func() {
	var builder *SinkStepBuilder

	BeforeEach(func() {
		builder = SinkStep("test-sink")
	})

	Describe("SinkStep", func() {
		It("should create a new sink step builder with type", func() {
			Expect(builder).ToNot(BeNil())
			Expect(builder.res).ToNot(BeNil())
			Expect(builder.res.Type).To(Equal("test-sink"))
		})
	})

	Describe("Configuration methods", func() {
		It("should set string values", func() {
			result := builder.SetString("endpoint", "http://example.com").Build()
			Expect(result.Config["endpoint"]).To(Equal("http://example.com"))
		})

		It("should set string array values", func() {
			result := builder.SetStrings("topics", "topic1", "topic2").Build()
			Expect(result.Config["topics"]).To(Equal([]string{"topic1", "topic2"}))
		})

		It("should set integer values", func() {
			result := builder.SetInt("batch_size", 100).Build()
			Expect(result.Config["batch_size"]).To(Equal(100))
		})

		It("should set boolean values", func() {
			result := builder.SetBool("compress", true).Build()
			Expect(result.Config["compress"]).To(Equal(true))
		})
	})

	Describe("Fluent API", func() {
		It("should support method chaining", func() {
			result := builder.
				SetString("host", "localhost").
				SetInt("port", 8080).
				SetBool("secure", true).
				SetStrings("allowed_methods", "GET", "POST")

			Expect(result).To(Equal(builder))

			sink := result.Build()
			Expect(sink.Config["host"]).To(Equal("localhost"))
			Expect(sink.Config["port"]).To(Equal(8080))
			Expect(sink.Config["secure"]).To(Equal(true))
			Expect(sink.Config["allowed_methods"]).To(Equal([]string{"GET", "POST"}))
		})
	})

	Describe("Edge cases", func() {
		It("should handle nil config initialization", func() {
			builder.res.Config = nil
			builder.SetString("test", "value")
			Expect(builder.res.Config).ToNot(BeNil())
		})

		It("should overwrite existing values", func() {
			builder.SetString("key", "value1")
			builder.SetString("key", "value2")

			result := builder.Build()
			Expect(result.Config["key"]).To(Equal("value2"))
		})
	})
})
