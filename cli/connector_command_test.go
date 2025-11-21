package cli

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect/v2/model"
)

var _ = Describe("ConnectorCommand", func() {
	var (
		opts   *Options
		cmd    *connectorCommand
		mockCl *mockClient
		appCtx *AppContext
	)

	BeforeEach(func() {
		// Setup test environment
		opts = &Options{
			Timeout: 5 * time.Second,
		}

		appCtx, mockCl = newMockAppContext()

		cmd = &connectorCommand{
			opts:    opts,
			envVars: make(map[string]string),
		}
	})

	Describe("listConnectors", func() {
		It("should list connectors when they exist", func() {
			mockCl.connectors = []model.ConnectorSummary{
				{
					ConnectorId: "test-connector-1",
					Description: "Test Connector 1",
					RuntimeId:   "synadia",
					Instances: model.ConnectorSummaryInstances{
						Running: 1,
						Stopped: 0,
						Pending: 0,
					},
				},
				{
					ConnectorId: "test-connector-2",
					Description: "Test Connector 2",
					RuntimeId:   "wombat",
					Instances: model.ConnectorSummaryInstances{
						Running: 0,
						Stopped: 1,
						Pending: 0,
					},
				},
			}

			err := cmd.listConnectorsWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle empty connector list", func() {
			mockCl.connectors = []model.ConnectorSummary{}

			err := cmd.listConnectorsWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle errors from client", func() {
			mockCl.connectorError = fmt.Errorf("connection failed")

			err := cmd.listConnectorsWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("connection failed"))
		})
	})

	Describe("getConnector", func() {
		BeforeEach(func() {
			cmd.id = "test-connector"
		})

		It("should get and display a connector", func() {
			mockCl.connector = &model.Connector{
				ConnectorId: "test-connector",
				Description: "Test Connector",
				RuntimeId:   "synadia",
				Steps: model.Steps{
					Source: &model.SourceStep{
						Type: "generate",
					},
				},
			}

			err := cmd.getConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle connector not found", func() {
			mockCl.connector = nil
			mockCl.connectorError = fmt.Errorf("connector not found")

			err := cmd.getConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("connector not found"))
		})
	})

	Describe("removeConnector", func() {
		BeforeEach(func() {
			cmd.id = "test-connector"
		})

		It("should delete a connector", func() {
			err := cmd.removeConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.deleteCalled).To(BeTrue())
		})

		It("should handle delete errors", func() {
			mockCl.connectorError = fmt.Errorf("cannot delete running connector")

			err := cmd.removeConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot delete running connector"))
		})
	})

	Describe("connectorStatus", func() {
		BeforeEach(func() {
			cmd.id = "test-connector"
		})

		It("should get connector status", func() {
			mockCl.connectorStatus = &model.ConnectorStatus{
				Running: 1,
				Stopped: 0,
			}

			err := cmd.connectorStatusWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle status errors", func() {
			mockCl.connectorError = fmt.Errorf("status unavailable")

			err := cmd.connectorStatusWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("status unavailable"))
		})
	})

	Describe("startConnector", func() {
		BeforeEach(func() {
			cmd.id = "test-connector"
			cmd.replicas = 1
			cmd.noPull = false
			cmd.envVars = make(map[string]string)
			cmd.startTimeout = "1m"
		})

		It("should start a connector with default options", func() {
			mockCl.instances = []model.Instance{
				{
					Id:          "instance-1",
					ConnectorId: "test-connector",
				},
			}

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.startCalled).To(BeTrue())
			Expect(mockCl.startOptions).ToNot(BeNil())
			Expect(mockCl.startOptions.Pull).To(BeTrue())
			Expect(mockCl.startOptions.Replicas).To(Equal(1))
		})

		It("should handle no-pull option", func() {
			cmd.noPull = true

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.startOptions.Pull).To(BeFalse())
		})

		It("should handle environment variables", func() {
			cmd.envVars["KEY1"] = "value1"
			cmd.envVars["KEY2"] = "value2"

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.startOptions.EnvVars).To(HaveLen(2))
			Expect(mockCl.startOptions.EnvVars["KEY1"]).To(Equal("value1"))
			Expect(mockCl.startOptions.EnvVars["KEY2"]).To(Equal("value2"))
		})

		It("should handle placement tags", func() {
			cmd.placementTags = []string{"tag1", "tag2"}

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.startOptions.PlacementTags).To(Equal([]string{"tag1", "tag2"}))
		})

		It("should parse start timeout", func() {
			cmd.startTimeout = "30s"

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.startOptions.Timeout).To(Equal("30s"))
		})

		It("should handle invalid timeout", func() {
			cmd.startTimeout = "invalid"

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid timeout"))
		})

		It("should handle start errors", func() {
			mockCl.connectorError = fmt.Errorf("failed to start")

			err := cmd.startConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to start"))
		})
	})

	Describe("stopConnector", func() {
		BeforeEach(func() {
			cmd.id = "test-connector"
		})

		It("should stop a connector", func() {
			mockCl.instances = []model.Instance{
				{
					Id:          "instance-1",
					ConnectorId: "test-connector",
				},
			}

			err := cmd.stopConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.stopCalled).To(BeTrue())
		})

		It("should handle stop errors", func() {
			mockCl.connectorError = fmt.Errorf("failed to stop")

			err := cmd.stopConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to stop"))
		})
	})

	Describe("copyConnector", func() {
		BeforeEach(func() {
			cmd.id = "source-connector"
			cmd.targetId = "target-connector"
		})

		It("should copy a connector", func() {
			// First it gets the source connector
			mockCl.connector = &model.Connector{
				ConnectorId:    "source-connector",
				Description:    "Source Connector",
				RuntimeId:      "synadia",
				RuntimeVersion: "v1.0.0",
				Steps: model.Steps{
					Source: &model.SourceStep{
						Type: "generate",
					},
				},
			}

			err := cmd.copyConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.createCalled).To(BeTrue())
		})

		It("should handle source connector not found", func() {
			mockCl.connectorError = fmt.Errorf("connector not found")

			err := cmd.copyConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("connector not found"))
		})
	})

	Describe("reloadConnector", func() {
		BeforeEach(func() {
			cmd.id = "test-connector"
		})

		It("should reload a connector by stopping and starting", func() {
			// The reload command should stop then start
			mockCl.instances = []model.Instance{
				{
					Id:          "instance-1",
					ConnectorId: "test-connector",
				},
			}

			err := cmd.reloadConnectorWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(mockCl.stopCalled).To(BeTrue())
			Expect(mockCl.startCalled).To(BeTrue())
		})

		It("should handle reload errors during stop", func() {
			mockCl.connectorError = fmt.Errorf("failed to stop")

			err := cmd.reloadConnectorWithClient(appCtx)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to stop"))
		})
	})

	Describe("selectConnectorTemplate", func() {
		It("should return selected template", func() {
			// This method would normally prompt the user
			// For testing, we can't easily mock the survey package
			// So we'll skip interactive tests
		})
	})
})
