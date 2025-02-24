# Consumer Configuration

The consumer section describes how to connect to NATS and read messages from it.

## Example

```yaml
subject: "my.cool.>"
nats_config:
  url: nats://nats.demo.io:4222
```

## Fields
| Field       | Type                                       | Default | Required | Description                                                                                                                                               |
|-------------|--------------------------------------------|---------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `subject`   | string                                     |         | yes      | The NATS subject from which to consume messages. This may contain wildcards.                                                                              |
| `nats`      | [NatsConfig](./nats_config.md)             |         | yes      | The configuration for the NATS connection                                                                                                                 |
| `queue`     | string                                     |         | no       | The queue this consumer belongs to. This becomes important when multiple executions of the connector are in play,                                         |
| `jetstream` | [JetstreamConfig](./consumer_jetstream.md) |         | no       | The jetstream producer options. Jetstream (and message acks) is only being used if this property is configured. Set it to `{}` to enable acknowledgements |