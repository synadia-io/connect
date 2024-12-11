# Steps Definition

Steps define the actual work that the connector will perform. Depending on the kind of connector, the steps can vary.
For example, an inlet will have at least a source and a producer, while an outlet will have at least a consumer and a 
sink. Certain connectors may have additional steps like a transformer to modify the messages as they flow through the
connector.

Beware that Source-Target or Consumer-Producer connectors are not allowed, you always need to have either a source/consumer
or a producer/sink defined.

## Example
### Inlet
```yaml
steps:
  source:
    type: ""
    config: {}

  transformer: {}
  
  producer:
    subject: ""
    nats_config:
      url: ""
      auth_enabled: false
      jwt: ""
      seed: ""
      username: ""
      password: ""
```

### Outlet
```yaml
steps:
  consumer:
    subject: ""
    nats_config:
      url: ""
      auth_enabled: false
      jwt: ""
      seed: ""
      username: ""
      password: ""
      
  transformer: {}
  
  sink:
    type: ""
    config: {}
```

## Fields
| Field         | Type                            | Default | Required | Description                                                                              |
|---------------|---------------------------------|---------|----------|------------------------------------------------------------------------------------------|
| `source`      | [Source](./source.md)           |         | inlet    | The source configuration for the inlet.                                                  |
| `consumer`    | [Consumer](./consumer.md)       |         | outlet   | The consumer information describing how to connect to NATS and how to read the messages. |
| `transformer` | [Transformer](./transformer.md) |         | no       | An optional transformer to change the messages as they flow through the connector        |                                                                                                                                                                                                                                                                                    |
| `producer`    | [Producer](./producer.md)       |         | inlet    | The producer configuration for the inlet.                                                |
| `sink`        | [Sink](./sink.md)               |         | outlet   | The sink describes how data should be written to the external system.                    |