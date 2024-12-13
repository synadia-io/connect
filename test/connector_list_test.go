package test

import (
	"context"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect-node/model"
	cl "github.com/synadia-io/connect/client"
	model2 "github.com/synadia-io/connect/model"
	"time"
)

var _ = Describe("Listing connectors", func() {
	When("no connectors are available", func() {
		It("should call the callback once with an empty message", func() {
			invocationCount := 0
			err := cc.ListConnectors(cl.ConnectorFilter{}, func(info *cl.ConnectorInfo, hasMore bool) error {
				invocationCount++
				Expect(info).To(BeNil())
				Expect(hasMore).To(BeFalse())
				return nil
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(invocationCount).To(Equal(1))
		})
	})

	When("matching connectors are available", func() {
		connectors := map[string]*model.Connector{}

		BeforeEach(func() {
			var err error
			var c *model.Connector

			c, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
			connectors[c.Id] = c

			c, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
			connectors[c.Id] = c
		})

		AfterEach(func() {
			for _, c := range connectors {
				_, err := l.DeleteConnector(context.Background(), clientAccount, c.Id)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should call the callback for each connector and and empty message", func() {
			invocationCount := 0
			detected := map[string]*cl.ConnectorInfo{}
			err := cc.ListConnectors(cl.ConnectorFilter{}, func(info *cl.ConnectorInfo, hasMore bool) error {
				invocationCount++

				if info != nil {
					Expect(info.Config.Steps).ToNot(BeNil())
					detected[info.ConnectorId] = info
				}

				return nil
			}, cl.WithTimeout(1*time.Minute))

			Expect(err).ToNot(HaveOccurred())
			Expect(invocationCount).To(Equal(3))

			for _, c := range connectors {
				det, fnd := detected[c.Id]
				if !fnd {
					Fail("Connector not found")
				}

				Expect(det.ConnectorId).To(Equal(c.Id))
				Expect(det.Description).To(Equal(c.Description))
				Expect(det.Kind).To(Equal(string(c.Kind())))
			}
		})
	})

	When("no matching connectors are available", func() {
		connectors := map[string]*model.Connector{}

		BeforeEach(func() {
			var err error
			var c *model.Connector

			c, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
			connectors[c.Id] = c

			c, _, err = l.CreateConnector(context.Background(), clientAccount, uuid.NewString(), internalInletConfig())
			Expect(err).ToNot(HaveOccurred())
			connectors[c.Id] = c
		})

		AfterEach(func() {
			for _, c := range connectors {
				_, err := l.DeleteConnector(context.Background(), clientAccount, c.Id)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should call the callback once with an empty message", func() {
			invocationCount := 0

			filter := cl.ConnectorFilter{
				Kinds: []model2.ConnectorKind{model2.Outlet},
			}

			err := cc.ListConnectors(filter, func(info *cl.ConnectorInfo, hasMore bool) error {
				invocationCount++

				return nil
			}, cl.WithTimeout(1*time.Minute))

			Expect(err).ToNot(HaveOccurred())
			Expect(invocationCount).To(Equal(1))
		})
	})
})
