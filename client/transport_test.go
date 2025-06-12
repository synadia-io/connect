package client

import (
	"time"

	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transport", func() {
	var (
		transport *Transport
		nc        *nats.Conn
	)

	BeforeEach(func() {
		// For unit tests, we'll create a transport without actual connection
		nc = &nats.Conn{}
		transport = &Transport{
			nc:      nc,
			trace:   false,
			account: "test-account",
		}
	})

	Describe("NewTransportForAccount", func() {
		It("should create a transport for specific account", func() {
			t := NewTransportForAccount(nc, "specific-account", false)
			Expect(t).ToNot(BeNil())
			Expect(t.Account()).To(Equal("specific-account"))
			Expect(t.trace).To(BeFalse())
		})

		It("should enable tracing when requested", func() {
			t := NewTransportForAccount(nc, "test-account", true)
			Expect(t.trace).To(BeTrue())
		})
	})

	Describe("Account", func() {
		It("should return the account name", func() {
			Expect(transport.Account()).To(Equal("test-account"))
		})
	})

	Describe("Request Options", func() {
		It("should provide default request options", func() {
			opts := DefaultRequestOpts()
			Expect(opts.Timeout).To(Equal(5 * time.Second))
		})

		It("should allow custom timeout", func() {
			opts := DefaultRequestOpts()
			WithTimeout(10 * time.Second)(opts)
			Expect(opts.Timeout).To(Equal(10 * time.Second))
		})
	})

	// Integration tests would require a real NATS server
	// These are unit tests focusing on the Transport structure and methods
})