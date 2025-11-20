package cli

import (
	"github.com/synadia-io/connect/v2/spec"
	"github.com/synadia-io/connect/v2/spec/builders"
)

var templates = []spec.ConnectorSpec{
	builders.Connector().
		Description("Inlet :: Send generated data to Core NATS").
		RuntimeId("wombat").
		Steps(builders.Steps().
			Source(builders.SourceStep("generate").
				SetString("mapping", "root.message = \"Hello World\"")).
			Producer(builders.ProducerStep(builders.NatsConfig("nats://demo.nats.io:4222")).
				Core(builders.ProducerStepCore("foo.bar")))).
		Build(),

	builders.Connector().
		Description("Outlet :: Send messages from Core NATS to MongoDB").
		RuntimeId("wombat").
		Steps(builders.Steps().
			Consumer(builders.ConsumerStep(builders.NatsConfig("nats://demo.nats.io:4222")).
				Core(builders.ConsumerStepCore("foo.bar"))).
			Sink(builders.SinkStep("mongodb").
				SetString("url", "mongodb+srv://your-mongo-server/?retryWrites=true").
				SetString("database", "my-db").
				SetString("collection", "my-collection").
				SetString("operation", "insert-one").
				SetString("document_map", "root = this"))).
		Build(),
}
