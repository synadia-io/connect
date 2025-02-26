# Consumer Jetstream Configuration
Additional options may be provided when consuming from a jetstream stream

## Fields
| Field             | Type   | Default | Required | Description                                                                                 |
|-------------------|--------|---------|----------|---------------------------------------------------------------------------------------------|
| `deliver_policy`  | string | `all`   | no       | The delivery policy for the consumer. Options are `all`, `last`, `new`                      |
| `max_ack_pending` | int    | 1       | no       | The maximum number of acks that can be pending before the consumer stops consuming messages |
| `max_ack_wait`    | string | `30s`   | no       | The maximum time to wait for an ack before retrying                                         |
| `durable`         | string |         | no       | The durable name for the consumer. If set, the consumer will be durable                     |