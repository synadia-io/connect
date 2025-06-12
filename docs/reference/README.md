# Synadia Connect Reference Documentation

This directory contains detailed reference documentation for Synadia Connect components and configurations.

## Component Reference

### Core Components

- [Connector](connector.md) - Complete connector configuration reference
- [Steps](steps.md) - Understanding connector steps and their configuration

### Data Flow Components

- [Source](source.md) - Components that read data from external systems
- [Consumer](consumer.md) - Components that consume data from NATS
- [Producer](producer.md) - Components that produce data to external systems
- [Transformer](transformer.md) - Components that transform data in transit

### Transformers

- [Composite Transformer](transformers/composite.md) - Combine multiple transformers
- [Mapping Transformer](transformers/mapping.md) - Map and transform data fields
- [Service Transformer](transformers/service.md) - Transform data using external services

### Configuration

- [NATS Configuration](nats_config.md) - NATS connection and authentication settings

## Component Categories

### Sources (Inlets)

Sources read data from external systems and publish to NATS:

- **HTTP Server** - Receive webhooks and HTTP requests
- **File** - Read from local or remote files
- **Database** - Query databases (PostgreSQL, MySQL, etc.)
- **Message Queues** - Consume from Kafka, RabbitMQ, etc.
- **APIs** - Poll REST APIs or GraphQL endpoints

### Consumers (Outlets)

Consumers read from NATS and process messages:

- **NATS Core** - Subscribe to NATS subjects
- **NATS JetStream** - Consume from JetStream streams
- **NATS Key/Value** - Watch KV bucket changes

### Producers (Outlets)

Producers write data to external systems:

- **HTTP Client** - Send HTTP requests
- **File** - Write to files
- **Database** - Insert/update database records
- **Message Queues** - Publish to external queues
- **Cloud Services** - Send to S3, GCS, Azure Blob, etc.

### Sinks (Inlets)

Sinks publish data to NATS:

- **NATS Core** - Publish to NATS subjects
- **NATS JetStream** - Publish to JetStream streams
- **NATS Key/Value** - Update KV buckets

## Configuration Examples

### Basic Inlet (HTTP → NATS)

```yaml
name: webhook-inlet
description: Receive webhooks and publish to NATS
runtime: synadia
config:
  source:
    type: http_server
    config:
      address: "0.0.0.0:8080"
      path: "/webhook"
  sink:
    type: nats_jetstream
    config:
      subject: "webhooks.received"
      stream: "WEBHOOKS"
```

### Basic Outlet (NATS → Database)

```yaml
name: events-to-postgres
description: Store NATS events in PostgreSQL
runtime: synadia
config:
  consumer:
    type: nats_jetstream
    config:
      stream: "EVENTS"
      consumer: "postgres-writer"
  producer:
    type: postgresql
    config:
      connection_string: "${DATABASE_URL}"
      table: "events"
      columns:
        - name: "id"
          source: "$.id"
        - name: "timestamp"
          source: "$.timestamp"
        - name: "data"
          source: "$"
```

### With Transformation

```yaml
name: transform-and-route
description: Transform messages and route to different subjects
runtime: synadia
config:
  source:
    type: nats_core
    config:
      subject: "raw.events"
  transformer:
    type: mapping
    config:
      rules:
        - input: "$.type"
          output: "event_type"
        - input: "$.payload"
          output: "data"
          transform: "uppercase"
  sink:
    type: nats_core
    config:
      subject_mapping:
        field: "event_type"
        prefix: "processed."
```

## Best Practices

### 1. Error Handling

Always configure error handling for production connectors:

```yaml
error_handling:
  max_retries: 3
  retry_delay: "5s"
  dead_letter_subject: "errors.${connector_id}"
```

### 2. Resource Management

Set appropriate resource limits:

```yaml
resources:
  memory: "256Mi"
  cpu: "100m"
```

### 3. Monitoring

Enable metrics for observability:

```yaml
metrics:
  enabled: true
  port: 8080
  path: "/metrics"
```

### 4. Security

Use environment variables for sensitive data:

```yaml
config:
  producer:
    type: postgresql
    config:
      connection_string: "${DATABASE_URL}"
      password: "${DB_PASSWORD}"
```

## Troubleshooting

### Common Issues

1. **Connection Failures**
   - Check NATS connectivity
   - Verify credentials and permissions
   - Ensure network policies allow traffic

2. **Performance Issues**
   - Monitor resource usage
   - Adjust batch sizes
   - Check for backpressure

3. **Data Loss**
   - Enable JetStream for persistence
   - Configure appropriate acknowledgment modes
   - Implement proper error handling

### Debug Mode

Enable debug logging for troubleshooting:

```bash
connect connector start my-connector --log-level debug
```

## Additional Resources

- [Connector Specification Schema](../../spec/schemas/connector-spec.schema.json)
- [Getting Started Guide](../getting-started.md)
- [NATS Documentation](https://docs.nats.io)
- [Connect CLI Reference](../cli-reference.md)