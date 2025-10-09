package builders

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SourceStepBuilder", func() {
	var builder *SourceStepBuilder

	BeforeEach(func() {
		builder = SourceStep("test-source")
	})

	Describe("SourceStep", func() {
		It("should create a new source step builder with type", func() {
			Expect(builder).ToNot(BeNil())
			Expect(builder.res).ToNot(BeNil())
			Expect(builder.res.Type).To(Equal("test-source"))
		})
	})

	Describe("SetString", func() {
		It("should set a string value", func() {
			builder.SetString("key", "value")

			result := builder.Build()
			Expect(result.Config).ToNot(BeNil())
			Expect(result.Config["key"]).To(Equal("value"))
		})

		It("should initialize config map if nil", func() {
			builder.res.Config = nil
			builder.SetString("key", "value")

			Expect(builder.res.Config).ToNot(BeNil())
			Expect(builder.res.Config["key"]).To(Equal("value"))
		})

		It("should support method chaining", func() {
			result := builder.SetString("key1", "value1").SetString("key2", "value2")

			Expect(result).To(Equal(builder))
			config := result.Build().Config
			Expect(config["key1"]).To(Equal("value1"))
			Expect(config["key2"]).To(Equal("value2"))
		})
	})

	Describe("SetStrings", func() {
		It("should set a string array value", func() {
			builder.SetStrings("keys", "value1", "value2", "value3")

			result := builder.Build()
			Expect(result.Config["keys"]).To(Equal([]string{"value1", "value2", "value3"}))
		})

		It("should handle empty array", func() {
			builder.SetStrings("keys")

			result := builder.Build()
			// When called with no variadic arguments, Go creates a nil slice
			Expect(result.Config["keys"]).To(BeNil())
		})

		It("should initialize config map if nil", func() {
			builder.res.Config = nil
			builder.SetStrings("keys", "value")

			Expect(builder.res.Config).ToNot(BeNil())
		})
	})

	Describe("SetInt", func() {
		It("should set an integer value", func() {
			builder.SetInt("count", 42)

			result := builder.Build()
			Expect(result.Config["count"]).To(Equal(42))
		})

		It("should handle negative values", func() {
			builder.SetInt("offset", -10)

			result := builder.Build()
			Expect(result.Config["offset"]).To(Equal(-10))
		})

		It("should initialize config map if nil", func() {
			builder.res.Config = nil
			builder.SetInt("count", 1)

			Expect(builder.res.Config).ToNot(BeNil())
		})
	})

	Describe("SetBool", func() {
		It("should set a boolean value to true", func() {
			builder.SetBool("enabled", true)

			result := builder.Build()
			Expect(result.Config["enabled"]).To(Equal(true))
		})

		It("should set a boolean value to false", func() {
			builder.SetBool("disabled", false)

			result := builder.Build()
			Expect(result.Config["disabled"]).To(Equal(false))
		})

		It("should initialize config map if nil", func() {
			builder.res.Config = nil
			builder.SetBool("flag", true)

			Expect(builder.res.Config).ToNot(BeNil())
		})
	})

	Describe("Build", func() {
		It("should return the built source step", func() {
			result := builder.
				SetString("url", "http://example.com").
				SetInt("timeout", 30).
				SetBool("verify_ssl", true).
				SetStrings("headers", "Content-Type: application/json").
				Build()

			Expect(result.Type).To(Equal("test-source"))
			Expect(result.Config["url"]).To(Equal("http://example.com"))
			Expect(result.Config["timeout"]).To(Equal(30))
			Expect(result.Config["verify_ssl"]).To(Equal(true))
			Expect(result.Config["headers"]).To(Equal([]string{"Content-Type: application/json"}))
		})

		It("should build with empty config", func() {
			newBuilder := SourceStep("empty-source")
			result := newBuilder.Build()

			Expect(result.Type).To(Equal("empty-source"))
			Expect(result.Config).To(BeNil())
		})
	})

	Describe("Complex configurations", func() {
		It("should handle mixed configuration types", func() {
			result := builder.
				SetString("name", "test").
				SetInt("retries", 3).
				SetBool("async", true).
				SetStrings("tags", "tag1", "tag2").
				SetString("description", "A test source").
				Build()

			Expect(result.Config).To(HaveLen(5))
			Expect(result.Config["name"]).To(Equal("test"))
			Expect(result.Config["retries"]).To(Equal(3))
			Expect(result.Config["async"]).To(Equal(true))
			Expect(result.Config["tags"]).To(Equal([]string{"tag1", "tag2"}))
			Expect(result.Config["description"]).To(Equal("A test source"))
		})
	})
})
