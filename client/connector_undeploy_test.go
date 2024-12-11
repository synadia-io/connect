package client_test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/model"
	"time"
)

var _ = Describe("Undeploy a connector", func() {
	When("no connector exists with the given ID", func() {
		It("should return an error", func() {
			dr, err := cc.UndeployConnector("non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(dr).To(BeNil())
		})
	})

	When("a connector exists with the given ID", func() {
		var connector *model.Connector
		var deployRequests []*model.DeployRequest

		BeforeEach(func() {
			var err error

			connector, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())

			// -- add a deploy capture
			capt.CaptureDeploy(func(req model.DeployRequest, payload model.DeployRequestPayload, msgs chan model.InstanceEvent, logs chan model.InstanceLog, metrics chan model.InstanceMetric) {
				deployRequests = append(deployRequests, &req)

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
			})
		})

		AfterEach(func() {
			capt.StopCaptureDeploy()
		})

		Context("and has not been deployed", func() {
			It("should return an error", func() {
				dr, err := cc.UndeployConnector("non-existent-id")
				Expect(err).To(HaveOccurred())
				Expect(dr).To(BeNil())
			})
		})

		Context("and has been deployed", func() {
			var did string

			BeforeEach(func() {
				dr, err := cc.DeployConnector(connector.Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(dr).ToNot(BeNil())
				Expect(dr.Error).To(Equal(""))
				Expect(dr.InstanceErrors).To(BeEmpty())
				Expect(dr.Targets).To(HaveLen(1))

				did = dr.DeploymentId

				// -- wait for the deploy to complete
				Eventually(func() bool {
					d, err := cc.GetDeployment(connector.Id, did)
					if err != nil {
						return false
					}

					return d.Status.Running > 0
				}, 5*time.Second, 100*time.Millisecond).Should(BeTrue())
			})
			It("should undeploy the connector", func() {
				ur, err := cc.UndeployConnector(connector.Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(ur).ToNot(BeNil())
				Expect(ur.DeploymentId).To(Equal(did))
				Expect(ur.Error).To(Equal(""))
				Expect(ur.InstanceErrors).To(BeEmpty())
			})
		})
	})
})
