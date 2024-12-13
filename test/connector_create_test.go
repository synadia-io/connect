package test

import (
	"context"
	"github.com/nats-io/nuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/model"
)

var _ = Describe("Creating a connector", func() {
	When("no connector exists with the given ID", func() {
		It("should create the connector", func() {
			id := nuid.Next()
			cfg := publicInletConfig()

			c, err := cc.CreateConnector(id, cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(c).ToNot(BeNil())
			Expect(c.Id).To(Equal(id))
			Expect(c.ConnectorConfig).To(BeEquivalentTo(cfg))
		})
	})

	When("a connector exists with the given ID", func() {
		var connector *model.Connector

		BeforeEach(func() {
			var err error

			connector, _, err = l.CreateConnector(context.Background(), clientAccount, nuid.Next(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error", func() {
			cfg := publicInletConfig()

			c, err := cc.CreateConnector(connector.Id, cfg)
			Expect(err).To(HaveOccurred())
			Expect(c).To(BeNil())
		})
	})
})
