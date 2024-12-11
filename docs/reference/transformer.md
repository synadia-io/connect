# Transformer
A transformer allows you to modify the data flowing through the connector.

Only a single transformer can be configured per connector. However, the `composite` transformer allows
you to chain multiple transformers together.

Currently, the following transformers are available:

- [`composite`](./transformers/composite.md): Chain multiple transformers together.
- [`mapping`](./transformers/mapping.md): Use the scripting to transform messages.
- [`service`](./transformers/service.md): Use a NATS service to transform messages.