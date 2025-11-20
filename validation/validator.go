package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/synadia-io/connect/v2/spec"
	"gopkg.in/yaml.v3"
)

type Validator struct {
	schemasPath string
}

func NewValidator() *Validator {
	return &Validator{
		schemasPath: "spec/schemas",
	}
}

func (v *Validator) ValidateConnectorFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var specData spec.Spec
	if err := yaml.Unmarshal(data, &specData); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	if specData.Type != spec.SpecTypeConnector {
		return fmt.Errorf("invalid spec type: expected %s, got %s", spec.SpecTypeConnector, specData.Type)
	}

	// Validate the connector spec structure
	if specData.Spec == nil {
		return fmt.Errorf("spec field is required")
	}

	return v.validateConnectorSpec(specData.Spec)
}

func (v *Validator) validateConnectorSpec(specData interface{}) error {
	// Convert to JSON for validation
	jsonData, err := json.Marshal(specData)
	if err != nil {
		return fmt.Errorf("failed to marshal spec: %w", err)
	}

	var connectorSpec spec.ConnectorSpec
	if err := json.Unmarshal(jsonData, &connectorSpec); err != nil {
		return fmt.Errorf("invalid connector spec format: %w", err)
	}

	// Basic validation
	if connectorSpec.Description == "" {
		return fmt.Errorf("description is required")
	}

	if connectorSpec.RuntimeId == "" {
		return fmt.Errorf("runtimeId is required")
	}

	// Validate steps
	if err := v.validateSteps(connectorSpec.Steps); err != nil {
		return fmt.Errorf("invalid steps: %w", err)
	}

	return nil
}

func (v *Validator) validateSteps(steps spec.StepsSpec) error {
	hasSource := steps.Source != nil
	hasSink := steps.Sink != nil
	hasConsumer := steps.Consumer != nil
	hasProducer := steps.Producer != nil

	// Must have at least one source/consumer and one sink/producer
	if !hasSource && !hasConsumer {
		return fmt.Errorf("connector must have either a source or consumer step")
	}

	if !hasSink && !hasProducer {
		return fmt.Errorf("connector must have either a sink or producer step")
	}

	// Validate individual steps
	if steps.Source != nil {
		if err := v.validateSourceStep(*steps.Source); err != nil {
			return fmt.Errorf("invalid source step: %w", err)
		}
	}

	if steps.Sink != nil {
		if err := v.validateSinkStep(*steps.Sink); err != nil {
			return fmt.Errorf("invalid sink step: %w", err)
		}
	}

	if steps.Consumer != nil {
		if err := v.validateConsumerStep(*steps.Consumer); err != nil {
			return fmt.Errorf("invalid consumer step: %w", err)
		}
	}

	if steps.Producer != nil {
		if err := v.validateProducerStep(*steps.Producer); err != nil {
			return fmt.Errorf("invalid producer step: %w", err)
		}
	}

	return nil
}

func (v *Validator) validateSourceStep(source spec.SourceStepSpec) error {
	if source.Type == "" {
		return fmt.Errorf("source type is required")
	}
	return nil
}

func (v *Validator) validateSinkStep(sink spec.SinkStepSpec) error {
	if sink.Type == "" {
		return fmt.Errorf("sink type is required")
	}
	return nil
}

func (v *Validator) validateConsumerStep(consumer spec.ConsumerStepSpec) error {
	if consumer.Nats.Url == "" {
		return fmt.Errorf("consumer NATS URL is required")
	}

	hasCore := consumer.Core != nil
	hasStream := consumer.Stream != nil
	hasKv := consumer.Kv != nil

	if !hasCore && !hasStream && !hasKv {
		return fmt.Errorf("consumer must have either core, stream, or kv configuration")
	}

	if hasCore && consumer.Core.Subject == "" {
		return fmt.Errorf("consumer core subject is required")
	}

	if hasStream && consumer.Stream.Subject == "" {
		return fmt.Errorf("consumer stream subject is required")
	}

	if hasKv && consumer.Kv.Bucket == "" {
		return fmt.Errorf("consumer kv bucket is required")
	}

	return nil
}

func (v *Validator) validateProducerStep(producer spec.ProducerStepSpec) error {
	if producer.Nats.Url == "" {
		return fmt.Errorf("producer NATS URL is required")
	}

	hasCore := producer.Core != nil
	hasStream := producer.Stream != nil
	hasKv := producer.Kv != nil

	if !hasCore && !hasStream && !hasKv {
		return fmt.Errorf("producer must have either core, stream, or kv configuration")
	}

	if hasCore && producer.Core.Subject == "" {
		return fmt.Errorf("producer core subject is required")
	}

	if hasStream && producer.Stream.Subject == "" {
		return fmt.Errorf("producer stream subject is required")
	}

	if hasKv && producer.Kv.Bucket == "" {
		return fmt.Errorf("producer kv bucket is required")
	}

	return nil
}

func (v *Validator) ValidateFileExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	return nil
}

func (v *Validator) ValidateFileExtension(filePath string) error {
	ext := filepath.Ext(filePath)
	if ext != ".yaml" && ext != ".yml" && ext != ".json" {
		return fmt.Errorf("unsupported file extension: %s (supported: .yaml, .yml, .json)", ext)
	}
	return nil
}
