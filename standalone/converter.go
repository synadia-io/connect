package standalone

import (
	"fmt"
	"strings"

	"github.com/synadia-io/connect/model"
)

// ConfigConverter converts Synadia Connect steps to runtime-specific configuration
type ConfigConverter interface {
	ConvertSteps(steps model.Steps) (string, error)
	GetRuntimeArgs() []string
}

// WombatConverter converts Synadia Connect steps to Wombat configuration
type WombatConverter struct{}

func NewWombatConverter() *WombatConverter {
	return &WombatConverter{}
}

func (c *WombatConverter) ConvertSteps(steps model.Steps) (string, error) {
	var config strings.Builder
	
	// Build Wombat configuration
	if steps.Source != nil {
		sourceConfig, err := c.convertSource(*steps.Source)
		if err != nil {
			return "", fmt.Errorf("failed to convert source: %w", err)
		}
		config.WriteString("input:\n")
		config.WriteString(sourceConfig)
		config.WriteString("\n")
	}

	if steps.Consumer != nil {
		consumerConfig, err := c.convertConsumer(*steps.Consumer)
		if err != nil {
			return "", fmt.Errorf("failed to convert consumer: %w", err)
		}
		config.WriteString("input:\n")
		config.WriteString(consumerConfig)
		config.WriteString("\n")
	}

	// Add processors for transformers
	if steps.Transformer != nil {
		processorConfig, err := c.convertTransformer(*steps.Transformer)
		if err != nil {
			return "", fmt.Errorf("failed to convert transformer: %w", err)
		}
		config.WriteString("pipeline:\n")
		config.WriteString("  processors:\n")
		config.WriteString(processorConfig)
		config.WriteString("\n")
	}

	if steps.Sink != nil {
		sinkConfig, err := c.convertSink(*steps.Sink)
		if err != nil {
			return "", fmt.Errorf("failed to convert sink: %w", err)
		}
		config.WriteString("output:\n")
		config.WriteString(sinkConfig)
		config.WriteString("\n")
	}

	if steps.Producer != nil {
		producerConfig, err := c.convertProducer(*steps.Producer)
		if err != nil {
			return "", fmt.Errorf("failed to convert producer: %w", err)
		}
		config.WriteString("output:\n")
		config.WriteString(producerConfig)
		config.WriteString("\n")
	}

	return config.String(), nil
}

func (c *WombatConverter) GetRuntimeArgs() []string {
	// Runtime args are not needed since we pass config as base64 argument
	return []string{}
}

func (c *WombatConverter) convertSource(source model.SourceStep) (string, error) {
	switch source.Type {
	case "http":
		return c.convertHTTPSource(source)
	case "file":
		return c.convertFileSource(source)
	default:
		return "", fmt.Errorf("unsupported source type: %s", source.Type)
	}
}

func (c *WombatConverter) convertHTTPSource(source model.SourceStep) (string, error) {
	config := "  http_server:\n"
	
	if port, ok := source.Config["port"]; ok {
		config += fmt.Sprintf("    address: \"0.0.0.0:%v\"\n", port)
	} else {
		config += "    address: \"0.0.0.0:8080\"\n"
	}
	
	if path, ok := source.Config["path"]; ok {
		config += fmt.Sprintf("    path: \"%v\"\n", path)
	} else {
		config += "    path: \"/\"\n"
	}
	
	return config, nil
}

func (c *WombatConverter) convertFileSource(source model.SourceStep) (string, error) {
	config := "  file:\n"
	
	if path, ok := source.Config["path"]; ok {
		config += fmt.Sprintf("    paths: [\"%v/**/*\"]\n", path)
	} else {
		config += "    paths: [\"/data/**/*\"]\n"
	}
	
	return config, nil
}

func (c *WombatConverter) convertConsumer(consumer model.ConsumerStep) (string, error) {
	config := "  nats:\n"
	config += fmt.Sprintf("    urls: [\"%s\"]\n", consumer.Nats.Url)
	
	if consumer.Core != nil {
		config += fmt.Sprintf("    subject: \"%s\"\n", consumer.Core.Subject)
		if consumer.Core.Queue != nil {
			config += fmt.Sprintf("    queue: \"%s\"\n", *consumer.Core.Queue)
		}
	} else if consumer.Stream != nil {
		config += fmt.Sprintf("    subject: \"%s\"\n", consumer.Stream.Subject)
		config += "    stream:\n"
		config += "      enabled: true\n"
	} else if consumer.Kv != nil {
		config += fmt.Sprintf("    kv_bucket: \"%s\"\n", consumer.Kv.Bucket)
		if consumer.Kv.Key != "" {
			config += fmt.Sprintf("    kv_key: \"%s\"\n", consumer.Kv.Key)
		}
	}
	
	// Add auth if configured
	if consumer.Nats.AuthEnabled {
		if consumer.Nats.Jwt != nil && consumer.Nats.Seed != nil {
			config += "    auth:\n"
			config += fmt.Sprintf("      nkey_file: \"%s\"\n", *consumer.Nats.Seed)
			config += fmt.Sprintf("      user_jwt: \"%s\"\n", *consumer.Nats.Jwt)
		}
	}
	
	return config, nil
}

func (c *WombatConverter) convertSink(sink model.SinkStep) (string, error) {
	switch sink.Type {
	case "http":
		return c.convertHTTPSink(sink)
	case "file":
		return c.convertFileSink(sink)
	case "database":
		return c.convertDatabaseSink(sink)
	default:
		return "", fmt.Errorf("unsupported sink type: %s", sink.Type)
	}
}

func (c *WombatConverter) convertHTTPSink(sink model.SinkStep) (string, error) {
	config := "  http_client:\n"
	
	if url, ok := sink.Config["url"]; ok {
		config += fmt.Sprintf("    url: \"%v\"\n", url)
	}
	
	if method, ok := sink.Config["method"]; ok {
		config += fmt.Sprintf("    verb: \"%v\"\n", method)
	} else {
		config += "    verb: \"POST\"\n"
	}
	
	return config, nil
}

func (c *WombatConverter) convertFileSink(sink model.SinkStep) (string, error) {
	config := "  file:\n"
	
	if path, ok := sink.Config["path"]; ok {
		config += fmt.Sprintf("    path: \"%v\"\n", path)
	} else {
		config += "    path: \"/data/output.txt\"\n"
	}
	
	return config, nil
}

func (c *WombatConverter) convertDatabaseSink(sink model.SinkStep) (string, error) {
	config := "  sql_insert:\n"
	
	if driver, ok := sink.Config["driver"]; ok {
		config += fmt.Sprintf("    driver: \"%v\"\n", driver)
	}
	
	if dsn, ok := sink.Config["dsn"]; ok {
		config += fmt.Sprintf("    dsn: \"%v\"\n", dsn)
	}
	
	if table, ok := sink.Config["table"]; ok {
		config += fmt.Sprintf("    table: \"%v\"\n", table)
	}
	
	return config, nil
}

func (c *WombatConverter) convertProducer(producer model.ProducerStep) (string, error) {
	config := "  nats:\n"
	config += fmt.Sprintf("    urls: [\"%s\"]\n", producer.Nats.Url)
	
	if producer.Core != nil {
		config += fmt.Sprintf("    subject: \"%s\"\n", producer.Core.Subject)
	} else if producer.Stream != nil {
		config += fmt.Sprintf("    subject: \"%s\"\n", producer.Stream.Subject)
		config += "    stream:\n"
		config += "      enabled: true\n"
	} else if producer.Kv != nil {
		config += fmt.Sprintf("    kv_bucket: \"%s\"\n", producer.Kv.Bucket)
		if producer.Kv.Key != "" {
			config += fmt.Sprintf("    kv_key: \"%s\"\n", producer.Kv.Key)
		}
	}
	
	// Add auth if configured
	if producer.Nats.AuthEnabled {
		if producer.Nats.Jwt != nil && producer.Nats.Seed != nil {
			config += "    auth:\n"
			config += fmt.Sprintf("      nkey_file: \"%s\"\n", *producer.Nats.Seed)
			config += fmt.Sprintf("      user_jwt: \"%s\"\n", *producer.Nats.Jwt)
		}
	}
	
	return config, nil
}

func (c *WombatConverter) convertTransformer(transformer model.TransformerStep) (string, error) {
	if transformer.Mapping != nil {
		return c.convertMappingTransformer(*transformer.Mapping)
	}
	// Add other transformer types as needed
	return "", fmt.Errorf("unsupported transformer type")
}

func (c *WombatConverter) convertMappingTransformer(mapping model.MappingTransformerStep) (string, error) {
	config := "    - mapping: |\n"
	// Add proper indentation to the mapping source code
	lines := strings.Split(mapping.Sourcecode, "\n")
	for _, line := range lines {
		config += fmt.Sprintf("        %s\n", line)
	}
	return config, nil
}

// GetConverter returns the appropriate converter for a runtime
func GetConverter(runtimeID string) (ConfigConverter, error) {
	// Parse runtime ID to get base runtime (remove version)
	parts := strings.SplitN(runtimeID, ":", 2)
	baseRuntime := parts[0]
	
	switch baseRuntime {
	case "wombat":
		return NewWombatConverter(), nil
	default:
		return nil, fmt.Errorf("no converter available for runtime: %s", baseRuntime)
	}
}