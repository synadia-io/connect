package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/synadia-io/connect/v2/spec"
	"github.com/synadia-io/connect/v2/spec/builders"
	"gopkg.in/yaml.v3"
)

func (c *standaloneCommand) loadConnectorSpec(filePath string) (*spec.ConnectorSpec, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var specData spec.Spec
	if err := yaml.Unmarshal(data, &specData); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if specData.Type != spec.SpecTypeConnector {
		return nil, fmt.Errorf("invalid spec type: expected %s, got %s", spec.SpecTypeConnector, specData.Type)
	}

	var connectorSpec spec.ConnectorSpec
	if err := mapstructure.Decode(specData.Spec, &connectorSpec); err != nil {
		return nil, fmt.Errorf("failed to decode connector spec: %w", err)
	}

	return &connectorSpec, nil
}

func (c *standaloneCommand) selectTemplate() (*spec.ConnectorSpec, error) {
	if c.templateName != "" {
		// Find template by name
		for _, template := range standaloneTemplates {
			if strings.Contains(strings.ToLower(template.Description), strings.ToLower(c.templateName)) {
				return &template.ConnectorSpec, nil
			}
		}
		return nil, fmt.Errorf("template '%s' not found", c.templateName)
	}

	// If no template specified, use the first one as default
	if len(standaloneTemplates) == 0 {
		return nil, fmt.Errorf("no templates available")
	}

	// Use default template (first one)
	defaultTemplate := standaloneTemplates[0]
	fmt.Printf("Using default template: %s\n", defaultTemplate.Description)
	return &defaultTemplate.ConnectorSpec, nil
}

func (c *standaloneCommand) writeConnectorFile(connectorSpec *spec.ConnectorSpec, filePath string) error {
	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: %s", filePath)
	}

	// Create the full spec structure
	fullSpec := spec.Spec{
		Type: spec.SpecTypeConnector,
		Spec: connectorSpec,
	}

	data, err := yaml.Marshal(fullSpec)
	if err != nil {
		return fmt.Errorf("failed to marshal spec: %w", err)
	}

	// Add header comment
	header := `# Synadia Connect Connector Definition
# This file defines a data pipeline connector for Synadia Connect
# 
# To validate: connect standalone validate <file>
# To run:      connect standalone run <file>
#
`

	finalContent := header + string(data)

	if err := os.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// connectorTemplate represents a built-in template
type connectorTemplate struct {
	Description string
	spec.ConnectorSpec
}

// Built-in templates for standalone mode
var standaloneTemplates = []connectorTemplate{
	{
		Description: "generate-to-nats",
		ConnectorSpec: builders.Connector().
			Description("Generate test messages and send to NATS Core").
			RuntimeId("wombat").
			Steps(builders.Steps().
				Source(builders.SourceStep("generate").
					SetString("mapping", "root.message = \"Hello from Wombat\"").
					SetString("interval", "5s")).
				Producer(builders.ProducerStep(builders.NatsConfig("nats://localhost:4222")).
					Core(builders.ProducerStepCore("events.test")))).
			Build(),
	},
	{
		Description: "nats-to-http",
		ConnectorSpec: builders.Connector().
			Description("Consume from NATS Core and send HTTP requests").
			RuntimeId("wombat").
			Steps(builders.Steps().
				Consumer(builders.ConsumerStep(builders.NatsConfig("nats://localhost:4222")).
					Core(builders.ConsumerStepCore("events.>").Queue("workers"))).
				Sink(builders.SinkStep("http_client").
					SetString("url", "http://localhost:3000/webhook").
					SetString("verb", "POST"))).
			Build(),
	},
	{
		Description: "nats-to-stream",
		ConnectorSpec: builders.Connector().
			Description("Forward messages from NATS Core to JetStream").
			RuntimeId("wombat").
			Steps(builders.Steps().
				Consumer(builders.ConsumerStep(builders.NatsConfig("nats://localhost:4222")).
					Core(builders.ConsumerStepCore("events.>").Queue("processors"))).
				Producer(builders.ProducerStep(builders.NatsConfig("nats://localhost:4222")).
					Stream(builders.ProducerStepStream("stream.events")))).
			Build(),
	},
	{
		Description: "nats-to-mongodb",
		ConnectorSpec: builders.Connector().
			Description("Generate test data and write to MongoDB").
			RuntimeId("wombat").
			Steps(builders.Steps().
				Consumer(builders.ConsumerStep(builders.NatsConfig("nats://localhost:4222")).
					Core(builders.ConsumerStepCore("events.>").Queue("workers"))).
				Sink(builders.SinkStep("mongodb").
					SetString("url", "mongodb://localhost:27017").
					SetString("database", "testdb").
					SetString("collection", "events").
					SetString("operation", "insert-one"))).
			Build(),
	},
}
