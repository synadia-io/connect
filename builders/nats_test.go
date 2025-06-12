package builders

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NatsConfigBuilder", func() {
	var builder *NatsConfigBuilder

	BeforeEach(func() {
		builder = NatsConfig()
	})

	Describe("NatsConfig", func() {
		It("should create a new NATS config builder with default URL", func() {
			Expect(builder).ToNot(BeNil())
			Expect(builder.nats).ToNot(BeNil())
			Expect(builder.nats.Url).To(Equal(DefaultNatsUrl))
		})
	})

	Describe("Url", func() {
		It("should set a custom URL", func() {
			builder.Url("nats://custom.server:4222")
			
			result := builder.Build()
			Expect(result.Url).To(Equal("nats://custom.server:4222"))
		})

		It("should support method chaining", func() {
			result := builder.Url("nats://example.com:4222")
			Expect(result).To(Equal(builder))
		})

		It("should overwrite the default URL", func() {
			result := builder.Url("nats://override:4222").Build()
			Expect(result.Url).To(Equal("nats://override:4222"))
			Expect(result.Url).ToNot(Equal(DefaultNatsUrl))
		})
	})

	Describe("Auth", func() {
		It("should set JWT and seed authentication", func() {
			jwt := "test-jwt-token"
			seed := "test-seed-value"
			
			builder.Auth(jwt, seed)
			
			result := builder.Build()
			Expect(result.AuthEnabled).To(BeTrue())
			Expect(result.Jwt).ToNot(BeNil())
			Expect(*result.Jwt).To(Equal(jwt))
			Expect(result.Seed).ToNot(BeNil())
			Expect(*result.Seed).To(Equal(seed))
		})

		It("should support method chaining with auth", func() {
			result := builder.Auth("jwt", "seed")
			Expect(result).To(Equal(builder))
		})
	})

	Describe("Build", func() {
		It("should build with default configuration", func() {
			result := builder.Build()
			
			Expect(result.Url).To(Equal(DefaultNatsUrl))
			Expect(result.AuthEnabled).To(BeFalse())
			Expect(result.Jwt).To(BeNil())
			Expect(result.Seed).To(BeNil())
		})

		It("should build with custom URL and auth", func() {
			result := builder.
				Url("nats://secure.server:4222").
				Auth("my-jwt", "my-seed").
				Build()

			Expect(result.Url).To(Equal("nats://secure.server:4222"))
			Expect(result.AuthEnabled).To(BeTrue())
			Expect(*result.Jwt).To(Equal("my-jwt"))
			Expect(*result.Seed).To(Equal("my-seed"))
		})
	})

	Describe("Fluent API", func() {
		It("should support full configuration chain", func() {
			config := NatsConfig().
				Url("nats://cluster.example.com:4222").
				Auth("jwt-token", "seed-key").
				Build()

			Expect(config.Url).To(Equal("nats://cluster.example.com:4222"))
			Expect(config.AuthEnabled).To(BeTrue())
			Expect(*config.Jwt).To(Equal("jwt-token"))
			Expect(*config.Seed).To(Equal("seed-key"))
		})
	})
})