package client_test

import (
	"context"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
	"github.com/nats-io/nkeys"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/agent"
	"github.com/synadia-io/connect-node/agent/capture"
	"github.com/synadia-io/connect-node/control"
	"github.com/synadia-io/connect-node/control/logic"
	im "github.com/synadia-io/connect-node/model"
	"github.com/synadia-io/connect-node/storage"
	cl "github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/test"
	"testing"
)

var srv *server.Server
var nc *nats.Conn
var kv jetstream.KeyValue

var l *logic.Logic
var xkey nkeys.KeyPair
var svc micro.Service

var clientAccount string
var cc cl.Client

var capt *capture.Capture
var agnt *agent.Agent
var agntXkey nkeys.KeyPair

var watcher *control.Watcher

func TestControl(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		acc := test.Account("TEST_CONTROL_ACCOUNT")
		clientAccount = acc.Id

		// -- start the nats server
		srv = test.NewDecentralizedServer().
			WithAccount(acc).
			Run()

		// -- connect to nats using the control account
		nc = acc.Connect(srv)

		// -- create the bucket
		js, err := jetstream.New(nc)
		Expect(err).ToNot(HaveOccurred())
		kv, err = js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
			Bucket: "CONNECTORS",
		})
		Expect(err).ToNot(HaveOccurred())

		// -- create the metastore
		store, err := storage.NewJetStreamMetastore(js, "CONNECTORS")
		Expect(err).ToNot(HaveOccurred())

		// -- generate the service key
		xkey, err = nkeys.CreateCurveKeys()
		Expect(err).ToNot(HaveOccurred())

		l, err = logic.NewLogic(nc, store, xkey)
		Expect(err).ToNot(HaveOccurred())

		// -- run the control service
		svc, err = control.RunService(nc, l, "0.0.0",
			func(service micro.Service) {},
			func(service micro.Service, natsError *micro.NATSError) {
				t.Log(natsError.Error())
				t.Fail()
			},
		)
		Expect(err).ToNot(HaveOccurred())

		cc, err = cl.NewClient(nc, true)
		Expect(err).ToNot(HaveOccurred())

		// -- create the agent
		capt = capture.New()
		agntXkey, err = nkeys.CreateCurveKeys()
		Expect(err).ToNot(HaveOccurred())

		agnt, err = agent.New(
			agent.WithNatsConn(nc),
			agent.WithExecutor(capt),
			agent.WithKeyPair(agntXkey),
		)
		Expect(err).ToNot(HaveOccurred())
		err = agnt.Listen(context.Background())
		Expect(err).To(BeNil())
		t.Log("agent listening")

		watcher = control.NewWatcher(nc, l)
		Expect(watcher.Watch()).ToNot(HaveOccurred())
	})

	AfterSuite(func() {
		// -- stop the agent
		agnt.Close()

		// -- stop the service
		svc.Stop()

		// -- close the nats connection
		nc.Close()

		// -- stop the server
		srv.Shutdown()
	})

	AfterEach(func() {
		// -- remove all connectors from the metastore
		keys, err := kv.Keys(context.Background())

		for _, key := range keys {
			err = kv.Delete(context.Background(), key)
			Expect(err).ToNot(HaveOccurred())

			err = kv.Purge(context.Background(), key)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	RunSpecs(t, "Control Suite")
}

func IsEqualToInternalConnector(pc *model.Connector, ic *im.Connector) {
	if pc != nil {
		Expect(pc.Id).To(Equal(ic.Id))
		Expect(pc.Description).To(Equal(ic.Description))
		Expect(pc.Workload).To(Equal(ic.Workload))

		if pc.Metrics != nil {
			Expect(pc.Metrics.Port).To(Equal(ic.Metrics.Port))
			Expect(pc.Metrics.Path).To(Equal(ic.Metrics.Path))
		} else {
			Expect(ic.Metrics).To(BeNil())
		}

		if pc.Steps != nil {
			if pc.Steps.Source != nil {
				Expect(pc.Steps.Source.Type).To(Equal(ic.Steps.Source.Type))
				Expect(pc.Steps.Source.Config).To(Equal(ic.Steps.Source.Config))
			} else {
				Expect(ic.Steps.Source).To(BeNil())
			}

			if pc.Steps.Producer != nil {
				IsEqualToInternalNatsConfig(pc.Steps.Producer.NatsConfig, ic.Steps.Producer.NatsConfig)

				if pc.Steps.Producer.JetStream != nil {
					pc.Steps.Producer.JetStream.AckWait = ic.Steps.Producer.JetStream.AckWait
					pc.Steps.Producer.JetStream.MsgId = ic.Steps.Producer.JetStream.MsgId

					if pc.Steps.Producer.JetStream.Batching != nil {
						Expect(pc.Steps.Producer.JetStream.Batching.Count).To(Equal(ic.Steps.Producer.JetStream.Batching.Count))
						Expect(pc.Steps.Producer.JetStream.Batching.ByteSize).To(Equal(ic.Steps.Producer.JetStream.Batching.ByteSize))
					} else {
						Expect(ic.Steps.Producer.JetStream.Batching).To(BeNil())
					}
				} else {
					Expect(ic.Steps.Producer.JetStream).To(BeNil())
				}

				Expect(pc.Steps.Producer.Subject).To(Equal(ic.Steps.Producer.Subject))
				Expect(pc.Steps.Producer.Threads).To(Equal(ic.Steps.Producer.Threads))
			} else {
				Expect(ic.Steps.Producer).To(BeNil())
			}

			IsEqualToInternalTransformer(pc.Steps.Transformer, ic.Steps.Transformer)

			if pc.Steps.Consumer != nil {
				IsEqualToInternalNatsConfig(pc.Steps.Consumer.NatsConfig, ic.Steps.Consumer.NatsConfig)

				if pc.Steps.Consumer.JetStream != nil {
					Expect(pc.Steps.Consumer.JetStream.DeliverPolicy).To(Equal(ic.Steps.Consumer.JetStream.DeliverPolicy))
					Expect(pc.Steps.Consumer.JetStream.MaxAckPending).To(Equal(ic.Steps.Consumer.JetStream.MaxAckPending))
					Expect(pc.Steps.Consumer.JetStream.MaxAckWait).To(Equal(ic.Steps.Consumer.JetStream.MaxAckWait))
					Expect(pc.Steps.Consumer.JetStream.Bind).To(Equal(ic.Steps.Consumer.JetStream.Bind))
					Expect(pc.Steps.Consumer.JetStream.Durable).To(Equal(ic.Steps.Consumer.JetStream.Durable))
				} else {
					Expect(ic.Steps.Consumer.JetStream).To(BeNil())
				}
			} else {
				Expect(ic.Steps.Consumer).To(BeNil())
			}

			if pc.Steps.Sink != nil {
				Expect(pc.Steps.Sink.Type).To(Equal(ic.Steps.Sink.Type))
				Expect(pc.Steps.Sink.Config).To(Equal(ic.Steps.Sink.Config))
			} else {
				Expect(ic.Steps.Sink).To(BeNil())
			}

		} else {
			Expect(ic.Steps).To(BeNil())
		}
	} else {
		Expect(ic).To(BeNil())
	}
}

func IsEqualToInternalNatsConfig(pc model.NatsConfig, ic im.NatsConfig) {
	Expect(pc.Url).To(Equal(ic.Url))
	Expect(pc.AuthEnabled).To(Equal(ic.AuthEnabled))
	Expect(pc.Jwt).To(Equal(ic.Jwt))
	Expect(pc.Seed).To(Equal(ic.Seed))
	Expect(pc.Username).To(Equal(ic.Username))
	Expect(pc.Password).To(Equal(ic.Password))
}

func IsEqualToInternalTransformer(pt *model.Transformer, it *im.Transformer) {
	if pt != nil {
		if pt.Service != nil {
			IsEqualToInternalNatsConfig(pt.Service.NatsConfig, it.Service.NatsConfig)
			Expect(pt.Service.Timeout).To(Equal(it.Service.Timeout))
			Expect(pt.Service.Endpoint).To(Equal(it.Service.Endpoint))
		} else {
			Expect(it.Service).To(BeNil())
		}

		if pt.Mapping != nil {
			Expect(pt.Mapping.Sourcecode).To(Equal(it.Mapping.Sourcecode))
		} else {
			Expect(it.Mapping).To(BeNil())
		}

		if pt.Composite != nil {
			if pt.Composite.Sequential != nil {
				for i, tr := range pt.Composite.Sequential {
					IsEqualToInternalTransformer(&tr, &it.Composite.Sequential[i])
				}
			} else {
				Expect(it.Composite).To(BeNil())
			}
		} else {
			Expect(it.Composite).To(BeNil())
		}
	} else {
		Expect(it).To(BeNil())
	}
}
