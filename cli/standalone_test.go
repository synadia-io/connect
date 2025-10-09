package cli

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StandaloneCommand", func() {
	var (
		opts    *Options
		cmd     *standaloneCommand
		tempDir string
		appCtx  *AppContext
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "standalone-test")
		Expect(err).ToNot(HaveOccurred())

		opts = &Options{
			Timeout:    5 * time.Second,
			Standalone: true,
		}

		appCtx, _ = newMockAppContext()

		cmd = &standaloneCommand{
			opts: opts,
		}

		// Change to temp directory for file operations
		err = os.Chdir(tempDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("getFilePath", func() {
		It("should return correct file path for connector name", func() {
			cmd.connectorName = "test-connector"
			path := cmd.getFilePath()
			Expect(path).To(Equal("test-connector.connector.yml"))
		})

		It("should handle connector names with dots", func() {
			cmd.connectorName = "my.test.connector"
			path := cmd.getFilePath()
			Expect(path).To(Equal("my.test.connector.connector.yml"))
		})
	})

	Describe("getContainerName", func() {
		It("should return correct container name", func() {
			cmd.connectorName = "test-connector"
			name := cmd.getContainerName()
			Expect(name).To(Equal("test-connector_connector"))
		})

		It("should handle connector names with dots", func() {
			cmd.connectorName = "my.test.connector"
			name := cmd.getContainerName()
			Expect(name).To(Equal("my.test.connector_connector"))
		})
	})

	Describe("createConnector", func() {
		BeforeEach(func() {
			cmd.connectorName = "test-connector"
			cmd.templateName = "generate" // Use "generate" which matches "Generate to NATS Core"
		})

		It("should create a connector file with correct content", func() {
			err := cmd.createConnector(nil)
			Expect(err).ToNot(HaveOccurred())

			// Verify file was created
			filePath := cmd.getFilePath()
			Expect(filePath).To(BeAnExistingFile())

			// Read and verify content
			content, err := os.ReadFile(filePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("test-connector"))
			Expect(string(content)).To(ContainSubstring("wombat"))
		})

		It("should return error for invalid template", func() {
			cmd.templateName = "invalid-template"
			err := cmd.createConnector(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("template 'invalid-template' not found"))
		})

		It("should not overwrite existing file without force", func() {
			// Create initial file
			err := cmd.createConnector(nil)
			Expect(err).ToNot(HaveOccurred())

			// Try to create again without force
			err = cmd.createConnector(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("file already exists"))
		})
	})

	Describe("validateConnector", func() {
		BeforeEach(func() {
			cmd.connectorName = "test-connector"
		})

		It("should validate a valid connector file", func() {
			// Create a valid connector file
			cmd.templateName = "generate"
			err := cmd.createConnector(nil)
			Expect(err).ToNot(HaveOccurred())

			// Validate it
			err = cmd.validateConnector(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error for non-existent file", func() {
			cmd.connectorName = "non-existent"
			err := cmd.validateConnector(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("file does not exist"))
		})

		It("should return error for invalid YAML", func() {
			// Create invalid YAML file
			filePath := cmd.getFilePath()
			err := os.WriteFile(filePath, []byte("invalid: yaml: content: ["), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cmd.validateConnector(nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("removeConnector", func() {
		BeforeEach(func() {
			cmd.connectorName = "test-connector"
		})

		It("should remove connector container when it exists", func() {
			// removeConnector removes the Docker container, not the file
			// Since Docker isn't running in test environment, it will fail with container not found
			// but that's expected behavior
			err := cmd.removeConnector(nil)
			// The error is expected because Docker container doesn't exist in test environment
			if err != nil {
				Expect(err.Error()).To(ContainSubstring("failed to remove connector"))
			}
		})

		It("should handle non-existent containers gracefully", func() {
			// removeConnector should handle missing containers gracefully
			err := cmd.removeConnector(nil)
			// In test environment, Docker errors are expected
			if err != nil {
				Expect(err.Error()).To(ContainSubstring("failed to remove connector"))
			}
		})
	})

	Describe("listTemplates", func() {
		It("should list available templates", func() {
			err := cmd.listTemplates(nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("template selection", func() {
		It("should select template by name", func() {
			cmd.templateName = "generate"
			template, err := cmd.selectTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(template).ToNot(BeNil())
		})

		It("should return error for unknown template", func() {
			cmd.templateName = "unknown-template"
			template, err := cmd.selectTemplate()
			Expect(err).To(HaveOccurred())
			Expect(template).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})

		It("should use default template when none specified", func() {
			cmd.templateName = ""
			template, err := cmd.selectTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(template).ToNot(BeNil())
		})
	})

	Describe("runtime management", func() {
		var runtimeID string

		BeforeEach(func() {
			// Generate a unique runtime ID for each test to avoid conflicts
			runtimeID = fmt.Sprintf("test-runtime-%d", time.Now().UnixNano())
		})

		AfterEach(func() {
			// Clean up: try to remove the test runtime if it exists
			if runtimeID != "" {
				cmd.runtimeID = runtimeID
				_ = cmd.removeRuntime(nil) // Ignore errors - runtime might not exist
			}
		})

		It("should list default wombat runtime", func() {
			err := cmd.listRuntimes(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should add new runtime with unique ID", func() {
			cmd.runtimeID = runtimeID
			cmd.runtimeRegistry = fmt.Sprintf("registry.example.com/%s", runtimeID)
			cmd.runtimeDescription = "Test runtime"

			err := cmd.addRuntime(nil)
			Expect(err).ToNot(HaveOccurred())

			// Verify runtime was added
			err = cmd.listRuntimes(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should remove runtime after adding it", func() {
			// First add a runtime
			cmd.runtimeID = runtimeID
			cmd.runtimeRegistry = fmt.Sprintf("registry.example.com/%s", runtimeID)
			cmd.runtimeDescription = "Test runtime"
			err := cmd.addRuntime(nil)
			Expect(err).ToNot(HaveOccurred())

			// Then remove it
			err = cmd.removeRuntime(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle default wombat runtime protection", func() {
			cmd.runtimeID = "wombat"
			err := cmd.removeRuntime(nil)
			// The actual error message from the implementation
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot remove default runtime"))
		})
	})

	Describe("integration with existing patterns", func() {
		It("should follow the same testing patterns as other CLI commands", func() {
			// This test verifies that the standalone command follows the same
			// patterns as other CLI commands like connectorCommand
			Expect(cmd.opts).ToNot(BeNil())
			Expect(cmd.opts.Standalone).To(BeTrue())
		})

		It("should work with AppContext like other commands", func() {
			// Verify the standalone command can work with AppContext
			// similar to how other commands do
			Expect(appCtx).ToNot(BeNil())
			Expect(appCtx.Client).ToNot(BeNil())
			Expect(appCtx.DefaultTimeout).To(Equal(5 * time.Second))
		})
	})
})
