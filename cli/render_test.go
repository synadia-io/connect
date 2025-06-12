package cli

import (
	"github.com/synadia-io/connect/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render", func() {
	Describe("renderBoolean", func() {
		It("should render nil as empty string", func() {
			result := renderBoolean(nil)
			Expect(result).To(Equal(""))
		})

		It("should render true as green yes", func() {
			b := true
			result := renderBoolean(&b)
			// The color package adds ANSI escape codes
			Expect(result).To(ContainSubstring("yes"))
		})

		It("should render false as red no", func() {
			b := false
			result := renderBoolean(&b)
			// The color package adds ANSI escape codes
			Expect(result).To(ContainSubstring("no"))
		})
	})

	Describe("renderConnector", func() {
		It("should render connector details", func() {
			connector := model.Connector{
				ConnectorId: "test-connector",
				Description: "Test connector for unit tests",
				RuntimeId:   "test-runtime",
				Steps: model.Steps{
					Source: &model.SourceStep{
						Type: "generate",
					},
				},
			}

			result := renderConnector(connector)
			Expect(result).To(ContainSubstring("Connector: test-connector"))
			Expect(result).To(ContainSubstring("Test connector for unit tests"))
			Expect(result).To(ContainSubstring("test-runtime"))
			Expect(result).To(ContainSubstring("source:"))
			Expect(result).To(ContainSubstring("type: generate"))
		})

		It("should handle empty steps", func() {
			connector := model.Connector{
				ConnectorId: "empty-connector",
				Description: "Connector with no steps",
				RuntimeId:   "test-runtime",
				Steps:       model.Steps{},
			}

			result := renderConnector(connector)
			Expect(result).To(ContainSubstring("Connector: empty-connector"))
			Expect(result).To(ContainSubstring("Connector with no steps"))
		})
	})

	Describe("renderRuntime", func() {
		It("should render runtime details", func() {
			description := "Test runtime description"
			email := "test@example.com"
			url := "https://example.com"
			
			runtime := model.Runtime{
				Id:          "test-runtime",
				Label:       "Test Runtime",
				Description: &description,
				Author: model.RuntimeAuthor{
					Name:  "Test Author",
					Email: &email,
					Url:   &url,
				},
				Image: "test/runtime:latest",
				Metrics: &model.RuntimeMetrics{
					Port: 8080,
				},
			}

			result := renderRuntime(runtime)
			Expect(result).To(ContainSubstring("Test Runtime"))
			Expect(result).To(ContainSubstring("test-runtime"))
			Expect(result).To(ContainSubstring("Test runtime description"))
			Expect(result).To(ContainSubstring("Test Author"))
			Expect(result).To(ContainSubstring("test@example.com"))
			Expect(result).To(ContainSubstring("https://example.com"))
			Expect(result).To(ContainSubstring("test/runtime:latest"))
			Expect(result).To(ContainSubstring("yes")) // Metrics enabled
		})

		It("should handle nil optional fields", func() {
			runtime := model.Runtime{
				Id:    "minimal-runtime",
				Label: "Minimal Runtime",
				Author: model.RuntimeAuthor{
					Name: "Author Name",
				},
				Image:   "minimal/runtime:latest",
				Metrics: nil,
			}

			result := renderRuntime(runtime)
			Expect(result).To(ContainSubstring("Minimal Runtime"))
			Expect(result).To(ContainSubstring("minimal-runtime"))
			Expect(result).To(ContainSubstring("Author Name"))
			Expect(result).ToNot(ContainSubstring("Description"))
			Expect(result).ToNot(ContainSubstring("Email"))
			Expect(result).ToNot(ContainSubstring("URL"))
			Expect(result).To(ContainSubstring("no")) // Metrics disabled
		})
	})
})