package builders

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConsumerStepBuilder", func() {
	var (
		builder    *ConsumerStepBuilder
		natsConfig *NatsConfigBuilder
	)

	BeforeEach(func() {
		natsConfig = NatsConfig().Url("nats://localhost:4222")
		builder = ConsumerStep(natsConfig)
	})

	Describe("ConsumerStep", func() {
		It("should create a new consumer step builder with NATS config", func() {
			Expect(builder).ToNot(BeNil())
			Expect(builder.res).ToNot(BeNil())
			Expect(builder.res.Nats).ToNot(BeNil())
			Expect(builder.res.Nats.Url).To(Equal("nats://localhost:4222"))
		})
	})

	Describe("Core", func() {
		It("should add a core consumer configuration", func() {
			core := ConsumerStepCore("test.subject")
			builder.Core(core)
			
			result := builder.Build()
			Expect(result.Core).ToNot(BeNil())
			Expect(result.Core.Subject).To(Equal("test.subject"))
		})

		It("should support method chaining", func() {
			core := ConsumerStepCore("test.subject")
			result := builder.Core(core)
			
			Expect(result).To(Equal(builder))
		})

		It("should support queue groups", func() {
			core := ConsumerStepCore("test.>").Queue("test-queue")
			builder.Core(core)
			
			result := builder.Build()
			Expect(result.Core).ToNot(BeNil())
			Expect(result.Core.Subject).To(Equal("test.>"))
			Expect(result.Core.Queue).ToNot(BeNil())
			Expect(*result.Core.Queue).To(Equal("test-queue"))
		})
	})

	Describe("Kv", func() {
		It("should add a key-value consumer configuration", func() {
			kv := ConsumerStepKv("test-bucket", "test-key")
			builder.Kv(kv)
			
			result := builder.Build()
			Expect(result.Kv).ToNot(BeNil())
			Expect(result.Kv.Bucket).To(Equal("test-bucket"))
			Expect(result.Kv.Key).To(Equal("test-key"))
		})
	})

	Describe("Stream", func() {
		It("should add a stream consumer configuration", func() {
			stream := ConsumerStepStream("test.stream.subject")
			builder.Stream(stream)
			
			result := builder.Build()
			Expect(result.Stream).ToNot(BeNil())
			Expect(result.Stream.Subject).To(Equal("test.stream.subject"))
		})
	})

	Describe("Build", func() {
		It("should build a consumer with only NATS config", func() {
			result := builder.Build()
			
			Expect(result.Nats).ToNot(BeNil())
			Expect(result.Core).To(BeNil())
			Expect(result.Kv).To(BeNil())
			Expect(result.Stream).To(BeNil())
		})

		It("should build a complete consumer configuration", func() {
			result := builder.
				Core(ConsumerStepCore("test.>")).
				Build()

			Expect(result.Nats.Url).To(Equal("nats://localhost:4222"))
			Expect(result.Core).ToNot(BeNil())
			Expect(result.Core.Subject).To(Equal("test.>"))
		})
	})

	Describe("Multiple consumer types", func() {
		It("should allow setting multiple consumer types", func() {
			// Note: In practice, only one consumer type should be used at a time
			// This test verifies the builder allows it technically
			result := builder.
				Core(ConsumerStepCore("core.subject")).
				Kv(ConsumerStepKv("kv-bucket", "kv-key")).
				Stream(ConsumerStepStream("stream.subject")).
				Build()

			Expect(result.Core).ToNot(BeNil())
			Expect(result.Kv).ToNot(BeNil())
			Expect(result.Stream).ToNot(BeNil())
		})
	})
})