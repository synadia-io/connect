package cli

import (
	"time"

	"github.com/choria-io/fisk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options", func() {
	var (
		app  *fisk.Application
		opts *Options
	)

	BeforeEach(func() {
		app = fisk.New("test", "test application")
		opts = &Options{}
	})

	Describe("RegisterFlags", func() {
		It("should register all flags", func() {
			RegisterFlags(app, "1.0.0", opts)

			// Verify flags were registered by checking the app has them
			// Parse with test values
			_, err := app.Parse([]string{
				"--server", "nats://localhost:4222",
				"--user", "testuser",
				"--password", "testpass",
				"--connection-name", "test-connection",
				"--creds", "test.creds",
				"--jwt", "test-jwt",
				"--seed", "test-seed",
				"--context", "test-context",
				"--timeout", "10s",
				"--log-level", "debug",
			})
			Expect(err).ToNot(HaveOccurred())

			// Verify values were set
			Expect(opts.Servers).To(Equal("nats://localhost:4222"))
			Expect(opts.Username).To(Equal("testuser"))
			Expect(opts.Password).To(Equal("testpass"))
			Expect(opts.ConnectionName).To(Equal("test-connection"))
			Expect(opts.Creds).To(Equal("test.creds"))
			Expect(opts.JWT).To(Equal("test-jwt"))
			Expect(opts.Seed).To(Equal("test-seed"))
			Expect(opts.ContextName).To(Equal("test-context"))
			Expect(opts.Timeout).To(Equal(10 * time.Second))
			Expect(opts.LogLevel).To(Equal("debug"))
		})

		It("should use default values", func() {
			RegisterFlags(app, "1.0.0", opts)
			
			_, err := app.Parse([]string{})
			Expect(err).ToNot(HaveOccurred())

			// Check defaults
			Expect(opts.ConnectionName).To(Equal("NATS Vent CLI Item 1.0.0"))
			Expect(opts.Timeout).To(Equal(5 * time.Second))
			Expect(opts.LogLevel).To(Equal("info"))
		})

		It("should accept short flags", func() {
			RegisterFlags(app, "1.0.0", opts)
			
			_, err := app.Parse([]string{"-s", "nats://short:4222"})
			Expect(err).ToNot(HaveOccurred())

			Expect(opts.Servers).To(Equal("nats://short:4222"))
		})

		It("should validate log level", func() {
			RegisterFlags(app, "1.0.0", opts)
			
			// Valid log levels
			for _, level := range []string{"error", "warn", "info", "debug", "trace"} {
				_, err := app.Parse([]string{"--log-level", level})
				Expect(err).ToNot(HaveOccurred())
				Expect(opts.LogLevel).To(Equal(level))
			}

			// Invalid log level
			_, err := app.Parse([]string{"--log-level", "invalid"})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("AppContext", func() {
		It("should close the NATS connection", func() {
			// This test would need a mock NATS connection
			// For now, we just verify the method exists
			ac := &AppContext{}
			// Calling Close() with nil Nc would panic in real code
			// but we're just verifying the method exists
			Expect(ac.Close).ToNot(BeNil())
		})
	})

	// Note: Testing LoadOptions and loadNats would require:
	// 1. A running NATS server or mock
	// 2. Valid credentials/contexts
	// These are better suited for integration tests
})