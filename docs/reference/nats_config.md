# NATS Configuration
The NATS configuration describes how to connect to a NATS server.

## Example

```yaml
url: nats://nats.demo.io:4222
auth_enabled: true
jwt: "eyJhb..."
seed: "..."
```

## Fields
| Field          | Type    | Default | Description                               |
|----------------|---------|---------|-------------------------------------------|
| `url`          | string  |         | The URL of the NATS server to connect to. |
| `auth_enabled` | boolean | false   | Whether to enable authentication.         |
| `jwt`          | string  |         | The JWT token to use for authentication.  |
| `seed`         | string  |         | The seed to use for authentication.       |
| `username`     | string  |         | The username to use for authentication.   |
| `password`     | string  |         | The password to use for authentication.   |
```