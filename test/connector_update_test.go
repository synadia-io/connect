package test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect/model"
)

var _ = Describe("Update a connector", func() {
	When("no connector exists with the given ID", func() {
		It("should not return an error", func() {
			id := uuid.NewString()
			cfg := publicInletConfig()

			c, err := cc.UpdateConnector(id, cfg)
			Expect(err).To(HaveOccurred())
			Expect(c).To(BeNil())
		})
	})

	When("a connector exists with the given ID", func() {
		var connector *model.Connector

		BeforeEach(func() {
			var err error

			connector, err = cc.CreateConnector(uuid.NewString(), publicInletConfig())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should update only the fields that were provided", func() {
			cfg := model.ConnectorConfig{
				Steps: &model.Steps{
					Source: &model.Source{
						Type: "updated-source-type",
					},
				},
			}

			c, err := cc.UpdateConnector(connector.Id, cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(c).ToNot(BeNil())

			connector.Steps.Source.Type = "updated-source-type"
			Expect(c).To(BeEquivalentTo(connector))
		})
	})
})
