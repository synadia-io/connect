# Metrics Definition

The metrics section of a connector definition allows you to configure how metrics are exposed by the connector. These
metrics can be retrieved through NATS, but the agent needs to know how to get them.

A metrics section is optional, but if you want to expose metrics, you need to configure it.

Runtimes can expose prometheus metrics on a specific port and path. This section allows you to configure the port and path
where the metrics will be exposed.

## Example
```yaml
metrics:
  port: 4195
  path: /metrics
```

## Fields
| Field  | Type   | Default | Required | Description                                                 |
|--------|--------|---------|----------|-------------------------------------------------------------|
| `port` | int    |         | no       | The port on which prometheus metrics will be made available |
| `path` | string |         | no       | The path to the prometheus metrics                          |