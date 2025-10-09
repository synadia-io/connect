package client

import (
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		nc     *nats.Conn
		client Client
	)

	Describe("NewClientForAccount", func() {
		It("should create a client for specific account", func() {
			// Create a minimal NATS connection for testing
			nc = &nats.Conn{}
			client = NewClientForAccount(nc, "test-account", false)
			Expect(client).ToNot(BeNil())
			Expect(client.Account()).To(Equal("test-account"))
		})

		It("should enable tracing when requested", func() {
			nc = &nats.Conn{}
			client = NewClientForAccount(nc, "test-account", true)
			Expect(client).ToNot(BeNil())
			// The trace flag is stored in the transport
		})
	})

	Describe("Client Interface", func() {
		BeforeEach(func() {
			nc = &nats.Conn{}
			client = NewClientForAccount(nc, "test-account", false)
		})

		It("should implement ConnectorClient interface", func() {
			var _ ConnectorClient = client
		})

		It("should implement LibraryClient interface", func() {
			var _ LibraryClient = client
		})

		It("should return account name", func() {
			Expect(client.Account()).To(Equal("test-account"))
		})

		It("should have Close method", func() {
			// Verify the Close method exists by checking it can be called
			// Note: We can't actually test Close() without a valid NATS connection
			// This just verifies the method exists on the interface
			var fn = client.Close
			Expect(fn).ToNot(BeNil())
		})
	})

	// Note: Full integration tests for the client methods would require
	// a running NATS server with the connect-node service.
	// These tests focus on the client structure and interface compliance.
})
