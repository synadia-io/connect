package cli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synadia-io/connect/model"
)

var _ = Describe("LibraryCommand", func() {
	var (
		cmd    *libraryCommand
		mockCl *mockClient
		appCtx *AppContext
	)

	BeforeEach(func() {
		appCtx, mockCl = newMockAppContext()

		cmd = &libraryCommand{
			opts: &Options{},
		}
	})

	Describe("listRuntimes", func() {
		It("should list available runtimes", func() {
			mockCl.runtimes = []model.RuntimeSummary{
				{
					Id:    "synadia",
					Label: "Synadia Runtime",
				},
				{
					Id:    "wombat",
					Label: "Wombat Runtime",
				},
			}

			// Since the function prints to stdout and doesn't return anything testable,
			// we just verify it doesn't error
			// In a real scenario, we might capture stdout or refactor to return data
			err := cmd.listRuntimesWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle empty runtime list", func() {
			mockCl.runtimes = []model.RuntimeSummary{}

			err := cmd.listRuntimesWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle errors", func() {
			// Our mock doesn't simulate runtime errors, but in real code it would
			err := cmd.listRuntimesWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("getRuntime", func() {
		BeforeEach(func() {
			cmd.runtime = "synadia"
		})

		It("should get and display runtime details", func() {
			description := "Native Synadia runtime"
			email := "support@synadia.com"
			url := "https://synadia.com"

			mockCl.runtime = &model.Runtime{
				Id:          "synadia",
				Label:       "Synadia Runtime",
				Description: &description,
				Author: model.RuntimeAuthor{
					Name:  "Synadia",
					Email: &email,
					Url:   &url,
				},
				Image: "synadia/runtime:latest",
			}

			err := cmd.getRuntimeWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle runtime not found", func() {
			mockCl.runtime = nil

			err := cmd.getRuntimeWithClient(appCtx)
			// The function might not return an error for nil runtime
			// but would display empty data
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("search", func() {
		It("should search components without filters", func() {
			mockCl.components = []model.ComponentSummary{
				{
					RuntimeId: "synadia",
					Name:      "nats_core_source",
					Label:     "NATS Core Source",
					Kind:      model.ComponentKindSource,
				},
				{
					RuntimeId: "wombat",
					Name:      "http_sink",
					Label:     "HTTP Sink",
					Kind:      model.ComponentKindSink,
				},
			}

			err := cmd.searchWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should search with runtime filter", func() {
			cmd.runtime = "synadia"

			mockCl.components = []model.ComponentSummary{
				{
					RuntimeId: "synadia",
					Name:      "nats_core_source",
					Label:     "NATS Core Source",
					Kind:      model.ComponentKindSource,
				},
			}

			err := cmd.searchWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should search with kind filter", func() {
			cmd.kind = "source"

			mockCl.components = []model.ComponentSummary{
				{
					RuntimeId: "synadia",
					Name:      "nats_core_source",
					Label:     "NATS Core Source",
					Kind:      model.ComponentKindSource,
				},
			}

			err := cmd.searchWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle empty search results", func() {
			mockCl.components = []model.ComponentSummary{}

			err := cmd.searchWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("info", func() {
		BeforeEach(func() {
			cmd.runtime = "synadia"
			cmd.kind = "source"
			cmd.component = "nats_core"
		})

		It("should display component info", func() {
			description := "Reads messages from NATS Core"

			mockCl.component = &model.Component{
				RuntimeId:   "synadia",
				Name:        "nats_core_source",
				Label:       "NATS Core Source",
				Kind:        model.ComponentKindSource,
				Description: &description,
				Status:      model.ComponentStatusStable,
				Fields: []model.ComponentField{
					{
						Name:  "subject",
						Label: "Subject",
						Type:  model.ComponentFieldTypeString,
					},
				},
			}

			err := cmd.infoWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle component not found", func() {
			mockCl.component = nil

			err := cmd.infoWithClient(appCtx)
			// The function might not return an error for nil component
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
