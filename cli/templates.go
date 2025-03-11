package cli

import (
    "github.com/synadia-io/connect/spec"
    . "github.com/synadia-io/connect/spec/builders"
)

var templates = []spec.ConnectorSpec{
    Connector().
        Description("Inlet :: Send generated data to Core NATS").
        RuntimeId("wombat:main").
        Steps(Steps().
            Source(SourceStep("generate").
                SetString("mapping", "root.message = \"Hello World\"")).
            Producer(ProducerStep(NatsConfig("nats://demo.nats.io:4222")).
                Core(ProducerStepCore("foo.bar")))).
        Build(),

    Connector().
        Description("Outlet :: Send messages from Core NATS to MongoDB").
        RuntimeId("wombat:main").
        Steps(Steps().
            Consumer(ConsumerStep(NatsConfig("nats://demo.nats.io:4222")).
                Core(ConsumerStepCore("foo.bar"))).
            Sink(SinkStep("mongodb").
                SetString("url", "mongodb+srv://your-mongo-server/?retryWrites=true").
                SetString("database", "my-db").
                SetString("collection", "my-collection").
                SetString("operation", "insert-one").
                SetString("document_map", "root = this"))).
        Build(),
}
