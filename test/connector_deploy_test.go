package test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/model"
	"time"
)

var _ = Describe("Deploy a connector", func() {
	When("no connector exists with the given ID", func() {
		It("should return an error", func() {
			dr, err := cc.DeployConnector("non-existent-id")
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

				msgs <- model.InstanceEvent{
					InstanceIdentity: req.InstanceIdentity,
					AgentId:          agntId,
					Type:             model.StoppedEventType,
					Timestamp:        time.Now(),
					ExitCode:         0,
				}
			})
		})

		AfterEach(func() {
			capt.StopCaptureDeploy()
		})

		It("should deploy the connector", func() {
			dr, err := cc.DeployConnector(connector.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(dr).ToNot(BeNil())
			Expect(dr.Error).To(Equal(""))
			Expect(dr.InstanceErrors).To(BeEmpty())
			Expect(dr.Targets).To(HaveLen(1))
		})
	})
})
