# Connector Definition

A Connector is a component that reads data from an external system and writes it to NATS or reads data from NATS and 
writes it to an external system.

Connectors are made up of steps that define the actual work that the connector will perform. Depending on the kind of
connector, the steps can vary. For example, an inlet will have at least a source and a producer, while an outlet will have
at least a consumer and a sink. Certain connectors may have additional steps like a transformer to modify the messages as
they flow through the connector.

A Connector is defined in a YAML file and can be created using the `connect create` command.

## Example
```yaml
description: ""
workload: ghcr.io/synadia-io/connect-runtime-vanilla:latest
metrics: {}
steps: {}
```

## Fields
| Field         | Type                    | Default | Required | Description                                                           |
|---------------|-------------------------|---------|----------|-----------------------------------------------------------------------|
| `description` | string                  |         | no       | A description of the inlet to provide more context to users.          |
| `workload`    | string                  |         | yes      | A reference to the runtime workload to use when running the connector |
| `metrics`     | [Metrics](./metrics.md) |         | no       | The metrics configuration for the inlet.                              |
| `steps`       | [Steps](./steps.md)     |         | yes      | The steps describing the work the connector wil perform.              |