package runtime_test

import (
	"context"
	"encoding/base64"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/runtime"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Runtime", func() {
	Expect(nil)
	When("launching a runtime", func() {
		It("should parse the given configuration", func() {
			cfg := `ZGVzY3JpcHRpb246IEEgc3VtbWFyeSBvZiB3aGF0IHRoaXMgY29ubmVjdG9yIGRvZXMKd29ya2xvYWQ6IHRlc3QKbWV0cmljczoKICAgIHBvcnQ6IDQxOTUKICAgIHBhdGg6IC9tZXRyaWNzCnN0ZXBzOgogICAgc291cmNlOgogICAgICAgIHR5cGU6IGdlbmVyYXRlCiAgICAgICAgY29uZmlnOgogICAgICAgICAgICBtYXBwaW5nOiB8CiAgICAgICAgICAgICAgICByb290ID0gIkhlbGxvIFdvcmxkIgogICAgY29uc3VtZXI6IG51bGwKICAgIHRyYW5zZm9ybWVyOiBudWxsCiAgICBwcm9kdWNlcjoKICAgICAgICBuYXRzX2NvbmZpZzoKICAgICAgICAgICAgdXJsOiBuYXRzOi8vZGVtby5uYXRzLmlvOjQyMjIKICAgICAgICAgICAgYXV0aGVuYWJsZWQ6IGZhbHNlCiAgICAgICAgICAgIGp3dDogIiIKICAgICAgICAgICAgc2VlZDogIiIKICAgICAgICAgICAgdXNlcm5hbWU6ICIiCiAgICAgICAgICAgIHBhc3N3b3JkOiAiIgogICAgICAgIHN1YmplY3Q6IGNvbm5lY3QudGVzdAogICAgICAgIGpldHN0cmVhbTogbnVsbAogICAgc2luazogbnVsbAo=`

			rt := runtime.NewRuntime("SYSTEM", "deplId", "execId", slog.LevelDebug)
			err := rt.Launch(context.Background(), func(ctx context.Context, runtime *runtime.Runtime, cfg model.ConnectorConfig) error {
				Expect(cfg.Description).To(Equal("A summary of what this connector does"))
				Expect(cfg.Workload).To(Equal("test"))
				Expect(cfg.Metrics.Port).To(Equal(4195))
				Expect(cfg.Metrics.Path).To(Equal("/metrics"))
				Expect(cfg.Steps.Source.Type).To(Equal("generate"))
				Expect(cfg.Steps.Source.Config["mapping"]).To(Equal("root = \"Hello World\"\n"))
				Expect(cfg.Steps.Producer.NatsConfig.Url).To(Equal("nats://demo.nats.io:4222"))
				Expect(cfg.Steps.Producer.Subject).To(Equal("connect.test"))
				return nil
			}, cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should parse the marshalled configuration from the deploy command", func() {
			cfg := []byte(`description: cool stuff from the CLI
workload: ghcr.io/synadia-io/connect-runtime-vanilla:latest
metrics:
    port: 4195
    path: /metrics
steps:
    source:
        type: generate
        config:
            interval: 1s
            mapping: |-
                let seq = counter(max: 10)
                root.name = "human%d".format($seq)
                root.message = "hello human %d".format($seq)
    producer:
        nats_config:
            url: nats://demo.nats.io:4222
        subject: testing.hello.${!name}`)

			internalCfg := &model.ConnectorConfig{}
			err := yaml.Unmarshal(cfg, internalCfg)
			Expect(err).ToNot(HaveOccurred())

			config, err := yaml.Marshal(internalCfg)
			Expect(err).ToNot(HaveOccurred())

			cfg = []byte(base64.StdEncoding.EncodeToString(config))

			rt := runtime.NewRuntime("SYSTEM", "deplId", "execId", slog.LevelDebug)
			err = rt.Launch(context.Background(), func(ctx context.Context, runtime *runtime.Runtime, cfg model.ConnectorConfig) error {
				Expect(cfg.Description).To(Equal("cool stuff from the CLI"))
				Expect(cfg.Workload).To(Equal("ghcr.io/synadia-io/connect-runtime-vanilla:latest"))
				Expect(cfg.Metrics.Port).To(Equal(4195))
				Expect(cfg.Metrics.Path).To(Equal("/metrics"))
				Expect(cfg.Steps.Source.Type).To(Equal("generate"))
				Expect(cfg.Steps.Source.Config["mapping"]).To(Equal("let seq = counter(max: 10)\nroot.name = \"human%d\".format($seq)\nroot.message = \"hello human %d\".format($seq)"))
				Expect(cfg.Steps.Producer.NatsConfig.Url).To(Equal("nats://demo.nats.io:4222"))
				Expect(cfg.Steps.Producer.Subject).To(Equal("testing.hello.${!name}"))
				return nil
			}, string(cfg))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
