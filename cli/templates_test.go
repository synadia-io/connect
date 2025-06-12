package cli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templates", func() {
	It("should have exactly 2 templates", func() {
		Expect(templates).To(HaveLen(2))
	})

	Describe("First template (Inlet)", func() {
		It("should be a valid inlet connector", func() {
			inlet := templates[0]
			
			Expect(inlet.Description).To(Equal("Inlet :: Send generated data to Core NATS"))
			Expect(inlet.RuntimeId).To(Equal("wombat"))
			Expect(inlet.Steps).ToNot(BeNil())
			
			// Check source
			Expect(inlet.Steps.Source).ToNot(BeNil())
			Expect(inlet.Steps.Source.Type).To(Equal("generate"))
			Expect(inlet.Steps.Source.Config).ToNot(BeNil())
			Expect(inlet.Steps.Source.Config["mapping"]).To(Equal("root.message = \"Hello World\""))
			
			// Check producer
			Expect(inlet.Steps.Producer).ToNot(BeNil())
			Expect(inlet.Steps.Producer.Nats.Url).To(Equal("nats://demo.nats.io:4222"))
			Expect(inlet.Steps.Producer.Core).ToNot(BeNil())
			Expect(inlet.Steps.Producer.Core.Subject).To(Equal("foo.bar"))
			
			// Should not have consumer or sink
			Expect(inlet.Steps.Consumer).To(BeNil())
			Expect(inlet.Steps.Sink).To(BeNil())
		})
	})

	Describe("Second template (Outlet)", func() {
		It("should be a valid outlet connector", func() {
			outlet := templates[1]
			
			Expect(outlet.Description).To(Equal("Outlet :: Send messages from Core NATS to MongoDB"))
			Expect(outlet.RuntimeId).To(Equal("wombat"))
			Expect(outlet.Steps).ToNot(BeNil())
			
			// Check consumer
			Expect(outlet.Steps.Consumer).ToNot(BeNil())
			Expect(outlet.Steps.Consumer.Nats.Url).To(Equal("nats://demo.nats.io:4222"))
			Expect(outlet.Steps.Consumer.Core).ToNot(BeNil())
			Expect(outlet.Steps.Consumer.Core.Subject).To(Equal("foo.bar"))
			
			// Check sink
			Expect(outlet.Steps.Sink).ToNot(BeNil())
			Expect(outlet.Steps.Sink.Type).To(Equal("mongodb"))
			Expect(outlet.Steps.Sink.Config).ToNot(BeNil())
			Expect(outlet.Steps.Sink.Config["url"]).To(Equal("mongodb+srv://your-mongo-server/?retryWrites=true"))
			Expect(outlet.Steps.Sink.Config["database"]).To(Equal("my-db"))
			Expect(outlet.Steps.Sink.Config["collection"]).To(Equal("my-collection"))
			Expect(outlet.Steps.Sink.Config["operation"]).To(Equal("insert-one"))
			Expect(outlet.Steps.Sink.Config["document_map"]).To(Equal("root = this"))
			
			// Should not have source or producer
			Expect(outlet.Steps.Source).To(BeNil())
			Expect(outlet.Steps.Producer).To(BeNil())
		})
	})
})