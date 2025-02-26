package runtime_test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"strings"

	"github.com/synadia-io/connect/v2/model"
	"github.com/synadia-io/connect/v2/runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Runtime", func() {
	Expect(nil)
	When("launching a runtime", func() {
		It("should parse the given configuration", func() {
			cfg := `cHJvZHVjZXI6CiAgICBuYXRzOgogICAgICAgIHVybDogbmF0czovL2RlbW8ubmF0cy5pbzo0MjIyCiAgICBzdWJqZWN0OiB0ZXN0aW5nLmhlbGxvLiR7IW5hbWV9CnNvdXJjZToKICAgIGNvbmZpZzoKICAgICAgICBpbnRlcnZhbDogMXMKICAgICAgICBtYXBwaW5nOiB8LQogICAgICAgICAgICBsZXQgc2VxID0gY291bnRlcihtYXg6IDEwKQogICAgICAgICAgICByb290Lm5hbWUgPSAiaHVtYW4lZCIuZm9ybWF0KCRzZXEpCiAgICAgICAgICAgIHJvb3QubWVzc2FnZSA9ICJoZWxsbyBodW1hbiAlZCIuZm9ybWF0KCRzZXEpCiAgICB0eXBlOiBnZW5lcmF0ZQo=`

			rt := runtime.NewRuntime(slog.LevelDebug)
			err := rt.Launch(context.Background(), func(ctx context.Context, runtime *runtime.Runtime, steps model.Steps) error {
				Expect(steps.Source).ToNot(BeNil())
				Expect(steps.Source.Type).To(Equal("generate"))
				Expect(steps.Source.Config["mapping"]).To(Equal("let seq = counter(max: 10)\nroot.name = \"human%d\".format($seq)\nroot.message = \"hello human %d\".format($seq)"))

				Expect(steps.Producer).ToNot(BeNil())
				Expect(steps.Producer.Nats.Url).To(Equal("nats://demo.nats.io:4222"))
				Expect(steps.Producer.Subject).To(Equal("testing.hello.${!name}"))
				return nil
			}, cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should parse the marshalled configuration from the deploy command", func() {
			cfg := []byte(strings.TrimSpace(`
source:
    type: generate
    config:
        interval: 1s
        mapping: |-
            let seq = counter(max: 10)
            root.name = "human%d".format($seq)
            root.message = "hello human %d".format($seq)
producer:
    nats:
        url: nats://demo.nats.io:4222
    subject: testing.hello.${!name}`))

			steps := model.Steps{}
			err := yaml.Unmarshal(cfg, &steps)
			Expect(err).ToNot(HaveOccurred())

			config, err := yaml.Marshal(steps)
			Expect(err).ToNot(HaveOccurred())

			cfg = []byte(base64.StdEncoding.EncodeToString(config))
			slog.Info("config", "cfg", string(cfg))
			rt := runtime.NewRuntime(slog.LevelDebug)
			err = rt.Launch(context.Background(), func(ctx context.Context, runtime *runtime.Runtime, steps model.Steps) error {
				Expect(steps.Source).ToNot(BeNil())
				Expect(steps.Source.Type).To(Equal("generate"))
				Expect(steps.Source.Config["mapping"]).To(Equal("let seq = counter(max: 10)\nroot.name = \"human%d\".format($seq)\nroot.message = \"hello human %d\".format($seq)"))

				Expect(steps.Producer).ToNot(BeNil())
				Expect(steps.Producer.Nats.Url).To(Equal("nats://demo.nats.io:4222"))
				Expect(steps.Producer.Subject).To(Equal("testing.hello.${!name}"))
				return nil
			}, string(cfg))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
