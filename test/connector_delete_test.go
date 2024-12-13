package test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/synadia-io/connect-node/control/logic"
	"github.com/synadia-io/connect-node/model"
	"time"
)

var _ = Describe("Deleting a connector", func() {
	When("no connector exists with the given ID", func() {
		It("should not return an error", func() {
			id := uuid.NewString()

			deleted, err := cc.DeleteConnector(id)
			Expect(err).ToNot(HaveOccurred())
			Expect(deleted).To(BeFalse())
		})
	})

	When("a connector exists with the given ID", func() {
		var connector *model.Connector

		BeforeEach(func() {
			var err error

			connector, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("and the connector is not deployed", func() {
			It("should remove the connector", func() {
				deleted, err := cc.DeleteConnector(connector.Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(deleted).To(BeTrue())

				c, _, err := l.GetConnector(context.Background(), clientAccount, connector.Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(c).To(BeNil())
			})
		})

		Context("and the connector is deployed", func() {
			BeforeEach(func() {
				// -- add a deploy capture
				capt.CaptureDeploy(func(req model.DeployRequest, payload model.DeployRequestPayload, msgs chan model.InstanceEvent, logs chan model.InstanceLog, metrics chan model.InstanceMetric) {
					agntId, _ := agntXkey.PublicKey()

					msgs <- model.InstanceEvent{
						InstanceIdentity: req.InstanceIdentity,
						AgentId:          agntId,
						Type:             model.CreatedEventType,
						Timestamp:        time.Now(),
					}

					msgs <- model.InstanceEvent{
						InstanceIdentity: req.InstanceIdentity,
						AgentId:          agntId,
						Type:             model.RunningEventType,
						Timestamp:        time.Now(),
					}

					msgs <- model.InstanceEvent{
						InstanceIdentity: req.InstanceIdentity,
						AgentId:          agntId,
						Type:             model.StoppedEventType,
						Timestamp:        time.Now(),
						ExitCode:         0,
					}
				})

				dr, err := l.DeployConnector(context.Background(), clientAccount, connector.Id, &logic.DeployConfig{
					Replicas: 1,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(dr.Error).To(Equal(""))
				Expect(dr.InstanceErrors).To(BeEmpty())

				// -- wait for the deploy to complete
				Eventually(func() bool {
					d, _, err := l.GetDeployment(context.Background(), clientAccount, dr.Deployment.ConnectorId, dr.Deployment.DeploymentId)
					if err != nil {
						GinkgoLogr.Error(err, "Failed to get deployment")
						return false
					}

					return d.Status.Running > 0
				}).Should(BeTrue())
				log.Info().Msg("Deployment is running")
			})

			AfterEach(func() {
				capt.StopCaptureDeploy()
			})

			It("should not remove the connector", func() {
				deleted, err := cc.DeleteConnector(connector.Id)
				Expect(err).To(HaveOccurred())
				Expect(deleted).To(BeFalse())

				c, _, err := l.GetConnector(context.Background(), clientAccount, connector.Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(c).ToNot(BeNil())
			})
		})
	})
})
