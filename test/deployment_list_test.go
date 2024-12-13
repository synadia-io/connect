package test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/control/logic"
	"github.com/synadia-io/connect-node/model"
	cl "github.com/synadia-io/connect/client"
)

var _ = Describe("Listing deployments", func() {
	When("no deployments exists", func() {
		It("should return nil and no error", func() {
			err := cc.ListDeployments(cl.DeploymentFilter{}, func(item *cl.DeploymentInfo, hasMore bool) error {
				Expect(item).To(BeNil())
				Expect(hasMore).To(BeFalse())
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("deployments exist", Ordered, func() {
		var connectors []*model.Connector
		var drs []*model.DeployResult

		BeforeAll(func() {
			var err error
			var conn *model.Connector
			var dr *model.DeployResult

			// first connector has one deployment
			conn, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
			connectors = append(connectors, conn)

			dr, err = l.DeployConnector(context.Background(), clientAccount, conn.Id, &logic.DeployConfig{})
			Expect(err).ToNot(HaveOccurred())
			drs = append(drs, dr)

			// second connector has two deployments
			conn, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
			connectors = append(connectors, conn)

			dr, err = l.DeployConnector(context.Background(), clientAccount, conn.Id, &logic.DeployConfig{})
			Expect(err).ToNot(HaveOccurred())
			drs = append(drs, dr)

			dr, err = l.DeployConnector(context.Background(), clientAccount, conn.Id, &logic.DeployConfig{})
			Expect(err).ToNot(HaveOccurred())
			drs = append(drs, dr)
		})

		Context("with a filter defined", func() {
			It("should return only the connectors that match the filter", func() {
				var count int
				filter := cl.DeploymentFilter{ConnectorId: connectors[1].Id}
				err := cc.ListDeployments(filter, func(item *cl.DeploymentInfo, hasMore bool) error {
					if item != nil {
						count++
					}

					return nil
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))
			})
		})
	})
})
