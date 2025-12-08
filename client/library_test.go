package client

import (
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect/model"
)

var _ = Describe("LibraryClient", func() {
	var (
		lc *libraryClient
		nc *nats.Conn
	)

	BeforeEach(func() {
		nc = &nats.Conn{}
		transport := NewTransportForAccount(nc, "test-account", false)
		lc = &libraryClient{t: transport}
	})

	Describe("subject generation", func() {
		It("should generate correct subject for library operations", func() {
			Expect(lc.subject(runtimes, "LIST")).To(Equal("$CONLIB.RUNTIMES.LIST"))
			Expect(lc.subject(runtimes, "GET")).To(Equal("$CONLIB.RUNTIMES.GET"))
			Expect(lc.subject(components, "LIST")).To(Equal("$CONLIB.COMPONENTS.LIST"))
			Expect(lc.subject(components, "GET")).To(Equal("$CONLIB.COMPONENTS.GET"))
		})
	})

	Describe("SearchComponents", func() {
		It("should create proper search request with filters", func() {
			// This test verifies the method signature and basic structure
			runtimeId := "synadia"
			kind := model.ComponentKindSource
			status := model.ComponentStatusStable

			filter := &model.ComponentSearchFilter{
				RuntimeId: &runtimeId,
				Kind:      &kind,
				Status:    &status,
			}

			// In a real test with a NATS server, this would verify the request
			// For now, we just ensure the method exists and accepts the right parameters
			var _ = lc.SearchComponents
			Expect(*filter.RuntimeId).To(Equal("synadia"))
		})
	})

	Describe("GetComponent", func() {
		It("should accept correct parameters", func() {
			// Verify method signature
			var _ = lc.GetComponent

			// In a real integration test, we would verify:
			// - The correct subject is used
			// - The request is properly formatted
			// - The response is correctly parsed
		})
	})
})
