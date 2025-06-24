# Synadia Connect API Documentation

This document provides a comprehensive reference for the Synadia Connect Go SDK and API.

## Installation

```bash
go get github.com/synadia-io/connect
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/nats-io/nats.go"
    "github.com/synadia-io/connect/client"
    "github.com/synadia-io/connect/model"
)

func main() {
    // Connect to NATS
    nc, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()
    
    // Create Connect client
    connectClient, err := client.NewClient(nc, false)
    if err != nil {
        log.Fatal(err)
    }
    
    // List connectors
    connectors, err := connectClient.ListConnectors(5 * time.Second)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, connector := range connectors {
        log.Printf("Connector: %s - %s\n", connector.ConnectorId, connector.Description)
    }
}
```

## Client SDK

### Creating a Client

```go
import "github.com/synadia-io/connect/client"

// Create with NATS connection
nc, _ := nats.Connect("nats://localhost:4222")
client, err := client.NewClient(nc, false) // false = no trace logging

// With trace logging enabled
clientWithTrace, err := client.NewClient(nc, true)
```

### Connector Operations

#### List Connectors

```go
connectors, err := client.ListConnectors(timeout)
for _, conn := range connectors {
    fmt.Printf("ID: %s, Runtime: %s, Running: %d\n", 
        conn.ConnectorId, conn.RuntimeId, conn.Instances.Running)
}
```

#### Get Connector

```go
connector, err := client.GetConnector("my-connector", timeout)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Connector: %+v\n", connector)
```

#### Create Connector

```go
steps := model.Steps{
    Source: &model.SourceStep{
        Type: "http_server",
        Config: map[string]interface{}{
            "address": "0.0.0.0:8080",
            "path": "/webhook",
        },
    },
    Sink: &model.SinkStep{
        Type: "nats_jetstream",
        Config: map[string]interface{}{
            "subject": "webhooks.received",
        },
    },
}

connector, err := client.CreateConnector(
    "webhook-inlet",
    "Receives webhooks and publishes to NATS",
    "synadia",
    steps,
    timeout,
)
```

#### Update Connector

```go
patch := `{
    "description": "Updated description",
    "steps": {
        "source": {
            "config": {
                "path": "/new-webhook"
            }
        }
    }
}`

updated, err := client.PatchConnector("my-connector", patch, timeout)
```

#### Delete Connector

```go
err := client.DeleteConnector("my-connector", timeout)
```

### Instance Management

#### Start Connector

```go
options := &model.ConnectorStartOptions{
    Pull:     true,
    Replicas: 3,
    Timeout:  "2m",
    EnvVars: model.ConnectorStartOptionsEnvVars{
        "API_KEY": "secret-key",
        "DEBUG":   "true",
    },
    PlacementTags: []string{"region:us-east", "env:prod"},
}

instances, err := client.StartConnector("my-connector", options, timeout)
```

#### Stop Connector

```go
instances, err := client.StopConnector("my-connector", timeout)
```

#### List Instances

```go
instances, err := client.ListConnectorInstances("my-connector", timeout)
for _, inst := range instances {
    fmt.Printf("Instance: %s\n", inst.Id)
}
```

#### Get Connector Status

```go
status, err := client.GetConnectorStatus("my-connector", timeout)
fmt.Printf("Running: %d, Stopped: %d\n", status.Running, status.Stopped)
```

### Library Operations

#### List Runtimes

```go
runtimes, err := client.ListRuntimes(timeout)
for _, rt := range runtimes {
    fmt.Printf("Runtime: %s - %s\n", rt.Id, rt.Label)
}
```

#### Get Runtime Details

```go
runtime, err := client.GetRuntime("synadia", timeout)
fmt.Printf("Runtime: %s, Image: %s\n", runtime.Id, runtime.Image)
```

#### Search Components

```go
// Search with filters
filter := &model.ComponentSearchFilter{
    RuntimeId: &[]string{"synadia"}[0],
    Kind:      &[]model.ComponentKind{model.ComponentKindSource}[0],
    Status:    &[]model.ComponentStatus{model.ComponentStatusStable}[0],
}

components, err := client.SearchComponents(filter, timeout)
```

#### Get Component Details

```go
component, err := client.GetComponent(
    "synadia",
    model.ComponentKindSource,
    "http_server",
    timeout,
)
```

## Builders API

The builders package provides a fluent API for constructing connectors:

```go
import "github.com/synadia-io/connect/builders"

// Create an inlet connector
inlet := builders.NewInletBuilder("webhook-inlet").
    Description("Webhook receiver").
    Source(builders.NewSourceBuilder("http_server").
        Config("address", "0.0.0.0:8080").
        Config("path", "/webhook").
        Build()).
    Transform(builders.NewTransformerBuilder("mapping").
        Config("rules", []map[string]string{
            {"input": "$.id", "output": "event_id"},
            {"input": "$.data", "output": "payload"},
        }).
        Build()).
    Sink(builders.NewSinkBuilder("nats_jetstream").
        Config("subject", "webhooks.processed").
        Config("stream", "WEBHOOKS").
        Build()).
    Build()

// Create an outlet connector
outlet := builders.NewOutletBuilder("db-writer").
    Description("Write events to database").
    Consumer(builders.NewConsumerBuilder("nats_jetstream").
        Config("stream", "EVENTS").
        Config("consumer", "db-writer").
        Build()).
    Producer(builders.NewProducerBuilder("postgresql").
        Config("connection_string", "${DATABASE_URL}").
        Config("table", "events").
        Build()).
    Build()
```

## Model Types

### Connector

```go
type Connector struct {
    ConnectorId string `json:"connector_id"`
    Description string `json:"description"`
    RuntimeId   string `json:"runtime_id"`
    Steps       Steps  `json:"steps"`
}
```

### Steps

```go
type Steps struct {
    Source      *SourceStep      `json:"source,omitempty"`
    Consumer    *ConsumerStep    `json:"consumer,omitempty"`
    Transformer *TransformerStep `json:"transformer,omitempty"`
    Sink        *SinkStep        `json:"sink,omitempty"`
    Producer    *ProducerStep    `json:"producer,omitempty"`
}
```

### Component Types

```go
type ComponentKind string

const (
    ComponentKindSource      ComponentKind = "source"
    ComponentKindSink        ComponentKind = "sink"
    ComponentKindConsumer    ComponentKind = "consumer"
    ComponentKindProducer    ComponentKind = "producer"
    ComponentKindTransformer ComponentKind = "transformer"
    ComponentKindScanner     ComponentKind = "scanner"
)

type ComponentStatus string

const (
    ComponentStatusStable       ComponentStatus = "stable"
    ComponentStatusPreview      ComponentStatus = "preview"
    ComponentStatusExperimental ComponentStatus = "experimental"
    ComponentStatusDeprecated   ComponentStatus = "deprecated"
)
```

## Error Handling

The SDK uses standard Go error handling patterns:

```go
connector, err := client.GetConnector("my-connector", timeout)
if err != nil {
    // Check for specific error types
    if errors.Is(err, nats.ErrTimeout) {
        // Handle timeout
    } else if strings.Contains(err.Error(), "not found") {
        // Handle not found
    } else {
        // Handle other errors
    }
}
```

## Best Practices

### 1. Timeout Management

Always use appropriate timeouts:

```go
const defaultTimeout = 5 * time.Second
const longOperationTimeout = 30 * time.Second

// Use longer timeout for operations that might take time
instances, err := client.StartConnector(id, opts, longOperationTimeout)
```

### 2. Resource Cleanup

Always close connections:

```go
nc, err := nats.Connect(servers)
if err != nil {
    return err
}
defer nc.Close()

client, err := client.NewClient(nc, false)
if err != nil {
    return err
}
// No need to close client - it uses the NATS connection
```

### 3. Environment Variables

Use environment variables for configuration:

```go
import "os"

servers := os.Getenv("NATS_URL")
if servers == "" {
    servers = "nats://localhost:4222"
}
```

### 4. Error Context

Add context to errors:

```go
connector, err := client.GetConnector(id, timeout)
if err != nil {
    return fmt.Errorf("failed to get connector %s: %w", id, err)
}
```

## Examples

See the [examples](examples/) directory for complete working examples:

- [Basic Connector Creation](examples/create-connector.go)
- [Connector Management](examples/manage-connectors.go)
- [Component Discovery](examples/discover-components.go)
- [Builder API Usage](examples/using-builders.go)

## Testing

When testing code that uses the Connect SDK, use mock implementations:

```go
type mockClient struct {
    connectors []model.ConnectorSummary
    err        error
}

func (m *mockClient) ListConnectors(timeout time.Duration) ([]model.ConnectorSummary, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.connectors, nil
}
```

## Support

- [GitHub Issues](https://github.com/synadia-io/connect/issues)
- [NATS Slack #connectors](https://slack.nats.io)