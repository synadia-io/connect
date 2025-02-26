# Service Transformer
The service transformer allows you to use a NATS service to transform messages. The message is sent to NATS as a
request and the response is used as the transformed message.

## Example
```yaml
service:
    subject: "services.my-service"
    nats_config:
        url: nats://nats.demo.io:4222
```

The example above shows how to use the `service` transformer to call a NATS service at the subject `services.my-service`.

## Fields
| Field      | Type                            | Default | Required | Description                                |
|------------|---------------------------------|---------|----------|--------------------------------------------|
| `endpoint` | string                          |         | yes      | The subject of the service to call.        |
| `nats`     | [NatsConfig](../nats_config.md) |         | yes      | The configuration for the NATS connection. |
| `timeout`  | string                          | "5s"    | no       | The maximum time to wait for a response.   |