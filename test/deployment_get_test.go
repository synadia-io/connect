package test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/model"
)

var _ = Describe("Getting a deployment", func() {
	When("no connector exists with the given ID", func() {
		It("should return nil and no error", func() {
			c, err := cc.GetDeployment("non-existent-id", "non-existing-deployment")
			Expect(err).ToNot(HaveOccurred())
			Expect(c).To(BeNil())
		})
	})

	When("connector exists but no deployment exists with the given ID", func() {
		var connector *model.Connector

		BeforeEach(func() {
			var err error

			connector, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return nil and no error", func() {
			c, err := cc.GetDeployment(connector.Id, "non-existent-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(c).To(BeNil())
		})
	})

	When("both connector and deployment exist", func() {
		var connector *model.Connector
		var did string

		BeforeEach(func() {
			var err error

			connector, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())

			dr, err := cc.DeployConnector(connector.Id)
			Expect(err).ToNot(HaveOccurred())
			did = dr.DeploymentId
		})

		It("should return the connector", func() {
			c, err := cc.GetDeployment(connector.Id, did)
			Expect(err).ToNot(HaveOccurred())
			Expect(c).ToNot(BeNil())
		})
	})
})
