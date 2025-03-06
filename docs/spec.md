# Connector Specification
A connector is defined in YAML. This document details the structure of the connector
configuration.

## Example
```yaml
description: A simple hello world inlet
runtime_id: wombat:main
steps:
  producer:
    core:
      subject: testing.hello.output
    nats:
      url: nats://demo.nats.io:4222
    threads: 1
  source:
    config:
      subject: testing.hello.input
      urls:
        - nats://demo.nats.io:4222
    type: nats
```

This example shows a simple connector listening on the `testing.hello.input` subject
and publishing messages to the `testing.hello.output` subject. It does however give an overview
of the high-level structure of a connector.

The full schema of the connector spec can be found [here](../spec/schemas/connector-spec.schema.json).

## Top-level Fields
### Description
Once you have a lot of connectors, it will become difficult to understand what each of them are 
supposed to do. The `description` field allows you to give a short description of what the connector

### Runtime ID
A connector is always linked to a runtime. Depending on the runtime, different sources and sinks will be
available to use within the connector. The `runtime_id` field specifies this runtime in the form `<id>:<version>`.
If no version is provided, `main` is being used as the default.

### Steps
The `steps` field defines where and how messages should be processed. A connector can have a source and producer or
a consumer and a sink. This means that a connector can read from a source and write to NATS or read from NATS and 
write to a sink.

## Consumers and Producers
Consumers and producers represent the NATS side of a connector. A consumer reads messages from NATS while a producer 
writes messages to NATS. Both consumers and producers can be configured to either use core NATS functionality, 
a JetStream stream or a JetStream KV.

[Producer Configuration Spec](../spec/schemas/connector-steps-producer-model.schema.json)
[Consumer Configuration Spec](../spec/schemas/connector-steps-consumer-model.schema.json)