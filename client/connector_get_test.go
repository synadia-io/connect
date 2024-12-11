package client_test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/model"
)

var _ = Describe("Getting a connector", func() {
	When("no connector exists with the given ID", func() {
		It("should return nil and no error", func() {
			c, err := cc.GetConnector("non-existent-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(c).To(BeNil())
		})
	})

	When("a connector exists with the given ID", func() {
		var connector *model.Connector

		BeforeEach(func() {
			var err error

			connector, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())

		})

		It("should return the connector", func() {
			existing, err := cc.GetConnector(connector.Id)
			Expect(err).ToNot(HaveOccurred())
			IsEqualToInternalConnector(existing, connector)
		})
	})
})
