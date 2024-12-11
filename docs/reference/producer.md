# Producer Configuration

The producer configuration describes how to write data to NATS.

## Example

```yaml
subject: "my.cool.subject"
nats_config:
  url: nats://nats.demo.io:4222
```

## Fields
| Field         | Type                           | Default | Description                                                                                                                                               |
|---------------|--------------------------------|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `subject`     | string                         |         | The NATS subject on which to publish messages.                                                                                                            |
| `nats_config` | [NatsConfig](./nats_config.md) |         | The configuration for the NATS connection                                                                                                                 |
| `threads`     | int                            | 1       | The number of threads to use for publishing. More threads should make things faster. If ordering is important, keep this to 1.                            |
| `jetstream`   | object                         |         | The jetstream producer options. Jetstream (and message acks) is only being used if this property is configured. Set it to `{}` to enable acknowledgements |