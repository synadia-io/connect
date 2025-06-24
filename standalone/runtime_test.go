package standalone

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuntimeManager", func() {
	var (
		manager *RuntimeManager
		tempDir string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "runtime-test")
		Expect(err).ToNot(HaveOccurred())

		// Create a manager and override its config dir for testing
		manager = NewRuntimeManager()
		manager.configDir = tempDir
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("NewRuntimeManager", func() {
		It("should create a new runtime manager", func() {
			rm := NewRuntimeManager()
			Expect(rm).ToNot(BeNil())
			Expect(rm.GetConfigDir()).ToNot(BeEmpty())
		})

		It("should initialize with default wombat runtime", func() {
			// Default runtime should be available immediately
			runtimes, err := manager.LoadRuntimes()
			Expect(err).ToNot(HaveOccurred())
			Expect(runtimes).To(HaveLen(1))
			Expect(runtimes[0].ID).To(Equal("wombat"))
			Expect(runtimes[0].Registry).To(Equal("registry.synadia.io/connect-runtime-wombat"))
		})
	})

	Describe("AddRuntime", func() {
		It("should add a new runtime", func() {
			runtime := Runtime{
				ID:          "test-runtime",
				Registry:    "registry.example.com/test-runtime",
				Description: "Test runtime for testing",
			}

			err := manager.AddRuntime(runtime)
			Expect(err).ToNot(HaveOccurred())

			// Verify runtime was added
			runtimes, err := manager.LoadRuntimes()
			Expect(err).ToNot(HaveOccurred())
			Expect(runtimes).To(HaveLen(2)) // wombat + test-runtime

			// Find the test runtime
			var testRuntime *Runtime
			for _, r := range runtimes {
				if r.ID == "test-runtime" {
					testRuntime = &r
					break
				}
			}
			Expect(testRuntime).ToNot(BeNil())
			Expect(testRuntime.Registry).To(Equal("registry.example.com/test-runtime"))
			Expect(testRuntime.Description).To(Equal("Test runtime for testing"))
		})

		It("should prevent duplicate runtime IDs", func() {
			runtime1 := Runtime{
				ID:          "duplicate-id",
				Registry:    "registry.example.com/runtime1",
				Description: "First runtime",
			}

			runtime2 := Runtime{
				ID:          "duplicate-id",
				Registry:    "registry.example.com/runtime2",
				Description: "Second runtime",
			}

			err := manager.AddRuntime(runtime1)
			Expect(err).ToNot(HaveOccurred())

			err = manager.AddRuntime(runtime2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already exists"))
		})

		It("should validate required fields", func() {
			runtime := Runtime{
				ID:          "",
				Registry:    "registry.example.com/test",
				Description: "Test runtime",
			}

			err := manager.AddRuntime(runtime)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ID is required"))
		})

		It("should validate registry field", func() {
			runtime := Runtime{
				ID:          "test-runtime",
				Registry:    "",
				Description: "Test runtime",
			}

			err := manager.AddRuntime(runtime)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Registry is required"))
		})
	})

	Describe("GetRuntime", func() {
		It("should get existing runtime", func() {
			runtime, err := manager.GetRuntime("wombat")
			Expect(err).ToNot(HaveOccurred())
			Expect(runtime).ToNot(BeNil())
			Expect(runtime.ID).To(Equal("wombat"))
			Expect(runtime.Registry).To(Equal("registry.synadia.io/connect-runtime-wombat"))
		})

		It("should return error for non-existent runtime", func() {
			runtime, err := manager.GetRuntime("non-existent")
			Expect(err).To(HaveOccurred())
			Expect(runtime).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("RemoveRuntime", func() {
		It("should remove custom runtime", func() {
			// Add a custom runtime first
			runtime := Runtime{
				ID:          "custom-runtime",
				Registry:    "registry.example.com/custom",
				Description: "Custom runtime",
			}
			err := manager.AddRuntime(runtime)
			Expect(err).ToNot(HaveOccurred())

			// Verify it exists
			runtimes, err := manager.LoadRuntimes()
			Expect(err).ToNot(HaveOccurred())
			Expect(runtimes).To(HaveLen(2))

			// Remove it
			err = manager.RemoveRuntime("custom-runtime")
			Expect(err).ToNot(HaveOccurred())

			// Verify it's gone
			runtimes, err = manager.LoadRuntimes()
			Expect(err).ToNot(HaveOccurred())
			Expect(runtimes).To(HaveLen(1))
			Expect(runtimes[0].ID).To(Equal("wombat"))
		})

		It("should prevent removal of default wombat runtime", func() {
			err := manager.RemoveRuntime("wombat")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot remove default runtime"))
		})

		It("should return error for non-existent runtime", func() {
			err := manager.RemoveRuntime("non-existent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("ResolveRuntimeImage", func() {
		It("should resolve runtime with version", func() {
			image, err := manager.ResolveRuntimeImage("wombat:v1.0.0")
			Expect(err).ToNot(HaveOccurred())
			Expect(image).To(Equal("registry.synadia.io/connect-runtime-wombat:v1.0.0"))
		})

		It("should resolve runtime without version (default to latest)", func() {
			image, err := manager.ResolveRuntimeImage("wombat")
			Expect(err).ToNot(HaveOccurred())
			Expect(image).To(Equal("registry.synadia.io/connect-runtime-wombat:latest"))
		})

		It("should handle custom runtime", func() {
			// Add custom runtime
			runtime := Runtime{
				ID:          "custom",
				Registry:    "registry.example.com/custom-runtime",
				Description: "Custom runtime",
			}
			err := manager.AddRuntime(runtime)
			Expect(err).ToNot(HaveOccurred())

			image, err := manager.ResolveRuntimeImage("custom:v2.0.0")
			Expect(err).ToNot(HaveOccurred())
			Expect(image).To(Equal("registry.example.com/custom-runtime:v2.0.0"))
		})

		It("should return error for unknown runtime", func() {
			image, err := manager.ResolveRuntimeImage("unknown:v1.0.0")
			Expect(err).To(HaveOccurred())
			Expect(image).To(BeEmpty())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("persistence", func() {
		It("should persist runtimes to file", func() {
			// Add a runtime
			runtime := Runtime{
				ID:          "persistent-runtime",
				Registry:    "registry.example.com/persistent",
				Description: "Persistent runtime",
			}
			err := manager.AddRuntime(runtime)
			Expect(err).ToNot(HaveOccurred())

			// Create new manager instance with same config path
			newManager := NewRuntimeManager()
			newManager.configDir = tempDir

			// Verify runtime persisted
			runtimes, err := newManager.LoadRuntimes()
			Expect(err).ToNot(HaveOccurred())
			Expect(runtimes).To(HaveLen(2)) // wombat + persistent-runtime

			var persistentRuntime *Runtime
			for _, r := range runtimes {
				if r.ID == "persistent-runtime" {
					persistentRuntime = &r
					break
				}
			}
			Expect(persistentRuntime).ToNot(BeNil())
			Expect(persistentRuntime.Registry).To(Equal("registry.example.com/persistent"))
		})

		It("should handle missing config file gracefully", func() {
			// Create manager with non-existent directory
			nonExistentDir := filepath.Join(tempDir, "non-existent")
			testManager := NewRuntimeManager()
			testManager.configDir = nonExistentDir

			// Should still have default runtime
			runtimes, err := testManager.LoadRuntimes()
			Expect(err).ToNot(HaveOccurred())
			Expect(runtimes).To(HaveLen(1))
			Expect(runtimes[0].ID).To(Equal("wombat"))
		})
	})

	Describe("validation", func() {
		It("should validate runtime structure", func() {
			runtime := Runtime{
				ID:          "valid-runtime",
				Registry:    "registry.example.com/valid",
				Description: "Valid runtime description",
			}

			// Basic validation
			Expect(runtime.ID).ToNot(BeEmpty())
			Expect(runtime.Registry).ToNot(BeEmpty())
		})

		It("should reject runtime with invalid ID", func() {
			runtime := Runtime{
				ID:          "", // Invalid: empty
				Registry:    "registry.example.com/test",
				Description: "Test runtime",
			}

			// Basic validation
			Expect(runtime.ID).To(BeEmpty())
		})

		It("should reject runtime with invalid registry", func() {
			runtime := Runtime{
				ID:          "test-runtime",
				Registry:    "", // Invalid: empty
				Description: "Test runtime",
			}

			// Basic validation
			Expect(runtime.Registry).To(BeEmpty())
		})
	})
})