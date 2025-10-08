package cli

import (
	"time"

	"github.com/choria-io/fisk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegisterFlags", func() {
	var (
		app  *fisk.Application
		opts *Options
	)

	BeforeEach(func() {
		app = fisk.New("test", "test application")
		opts = &Options{}
	})

	It("should register all global flags", func() {
		RegisterFlags(app, "test-version", opts)

		// Check that flags are registered by looking at the model
		flags := app.Model().Flags
		flagNames := []string{}
		for _, flag := range flags {
			flagNames = append(flagNames, flag.Name)
		}

		// Check all expected flags (server not servers)
		Expect(flagNames).To(ContainElement("server"))
		Expect(flagNames).To(ContainElement("context"))
		Expect(flagNames).To(ContainElement("timeout"))
		Expect(flagNames).To(ContainElement("log-level"))
		// version and help are handled by fisk itself
	})

	It("should parse server flag correctly", func() {
		RegisterFlags(app, "test-version", opts)

		_, err := app.Parse([]string{"--server", "nats://localhost:4222"})
		Expect(err).To(BeNil())
		Expect(opts.Servers).To(Equal("nats://localhost:4222"))
	})

	It("should parse multiple flags correctly", func() {
		RegisterFlags(app, "test-version", opts)

		_, err := app.Parse([]string{
			"--server", "nats://localhost:4222",
			"--context", "test-context",
			"--timeout", "10s",
			"--log-level", "debug",
		})
		Expect(err).To(BeNil())
		Expect(opts.Servers).To(Equal("nats://localhost:4222"))
		Expect(opts.ContextName).To(Equal("test-context"))
		Expect(opts.Timeout).To(Equal(10 * time.Second))
		Expect(opts.LogLevel).To(Equal("debug"))
	})

	It("should use correct short flags", func() {
		RegisterFlags(app, "test-version", opts)

		_, err := app.Parse([]string{
			"-s", "nats://localhost:4222",
		})
		Expect(err).To(BeNil())
		Expect(opts.Servers).To(Equal("nats://localhost:4222"))
	})

	It("should validate log level enum", func() {
		RegisterFlags(app, "test-version", opts)

		// Valid log levels
		for _, level := range []string{"error", "warn", "info", "debug", "trace"} {
			_, err := app.Parse([]string{"--log-level", level})
			Expect(err).To(BeNil())
			Expect(opts.LogLevel).To(Equal(level))
		}

		// Invalid log level
		_, err := app.Parse([]string{"--log-level", "invalid"})
		Expect(err).To(HaveOccurred())
	})

	It("should handle default timeout", func() {
		RegisterFlags(app, "test-version", opts)

		_, err := app.Parse([]string{})
		Expect(err).To(BeNil())
		// Default timeout is 30s
		Expect(opts.Timeout).To(Equal(30 * time.Second))
	})

	It("should handle connection name default", func() {
		RegisterFlags(app, "test-version", opts)

		_, err := app.Parse([]string{})
		Expect(err).To(BeNil())
		Expect(opts.ConnectionName).To(Equal("NATS Vent CLI Item test-version"))
	})
})
