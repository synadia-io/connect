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
			cfg := `cHJvZHVjZXI6CiAgICBjb3JlOgogICAgICAgIHN1YmplY3Q6IHRlc3RpbmcuaGVsbG8uJHshbmFtZX0KICAgIG5hdHM6CiAgICAgICAgdXJsOiBuYXRzOi8vZGVtby5uYXRzLmlvOjQyMjIKc291cmNlOgogICAgY29uZmlnOgogICAgICAgIGludGVydmFsOiAxcwogICAgICAgIG1hcHBpbmc6IHwtCiAgICAgICAgICAgIGxldCBzZXEgPSBjb3VudGVyKG1heDogMTApCiAgICAgICAgICAgIHJvb3QubmFtZSA9ICJodW1hbiVkIi5mb3JtYXQoJHNlcSkKICAgICAgICAgICAgcm9vdC5tZXNzYWdlID0gImhlbGxvIGh1bWFuICVkIi5mb3JtYXQoJHNlcSkKICAgIHR5cGU6IGdlbmVyYXRlCg==`

			rt := runtime.NewRuntime()
			err := rt.Launch(context.Background(), func(ctx context.Context, runtime *runtime.Runtime, steps model.Steps) error {
				Expect(steps.Source).ToNot(BeNil())
				Expect(steps.Source.Type).To(Equal("generate"))
				Expect(steps.Source.Config["mapping"]).To(Equal("let seq = counter(max: 10)\nroot.name = \"human%d\".format($seq)\nroot.message = \"hello human %d\".format($seq)"))

				Expect(steps.Producer).ToNot(BeNil())
				Expect(steps.Producer.Nats.Url).To(Equal("nats://demo.nats.io:4222"))
				Expect(steps.Producer.Core).ToNot(BeNil())
				Expect(steps.Producer.Core.Subject).To(Equal("testing.hello.${!name}"))
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
    core:
        subject: testing.hello.${!name}`))

			steps := model.Steps{}
			err := yaml.Unmarshal(cfg, &steps)
			Expect(err).ToNot(HaveOccurred())

			config, err := yaml.Marshal(steps)
			Expect(err).ToNot(HaveOccurred())

			cfg = []byte(base64.StdEncoding.EncodeToString(config))
			slog.Info("config", "cfg", string(cfg))
			rt := runtime.NewRuntime()
			err = rt.Launch(context.Background(), func(ctx context.Context, runtime *runtime.Runtime, steps model.Steps) error {
				Expect(steps.Source).ToNot(BeNil())
				Expect(steps.Source.Type).To(Equal("generate"))
				Expect(steps.Source.Config["mapping"]).To(Equal("let seq = counter(max: 10)\nroot.name = \"human%d\".format($seq)\nroot.message = \"hello human %d\".format($seq)"))

				Expect(steps.Producer).ToNot(BeNil())
				Expect(steps.Producer.Nats.Url).To(Equal("nats://demo.nats.io:4222"))
				Expect(steps.Producer.Core).ToNot(BeNil())
				Expect(steps.Producer.Core.Subject).To(Equal("testing.hello.${!name}"))
				return nil
			}, string(cfg))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
