# Quick Start Example

Create and run your first connector in 5 minutes using standalone mode.

## Overview

This example creates a simple data generator that:
1. Generates test messages every 5 seconds
2. Publishes them to a NATS subject
3. Shows how to monitor and manage the connector

## Prerequisites

- Connect CLI with standalone mode
- Docker installed and running

## Step 1: Create the Connector

```shell
# Create a new connector from the generate template
connect standalone create my-first-connector --template generate

# This creates: my-first-connector.connector.yml
```

## Step 2: Examine the Configuration

```shell
cat my-first-connector.connector.yml
```

You'll see something like:

```yaml
# Synadia Connect Connector Definition
type: connector
spec:
  id: my-first-connector
  description: Connector: my-first-connector
  runtime_id: wombat
  steps:
    source:
      type: generate
      config:
        mapping: 'root.message = "Hello from Wombat"'
        interval: 5s
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: events.test
```

## Step 3: Validate the Configuration

```shell
connect standalone validate my-first-connector
```

Expected output:
```
Validating connector 'my-first-connector' (file: my-first-connector.connector.yml)
✓ Connector 'my-first-connector' is valid
```

## Step 4: Start NATS Server (Optional)

If you want to see messages flowing, start a NATS server:

```shell
# Install NATS CLI if you don't have it
# brew install nats-io/nats-tools/nats

# Start NATS server with JetStream
nats server --jetstream &

# Subscribe to see messages (in another terminal)
nats sub events.test
```

## Step 5: Run the Connector

```shell
# Run the connector and follow logs
connect standalone run my-first-connector --follow
```

Expected output:
```
Using runtime 'wombat' resolved to image 'registry.synadia.io/connect-runtime-wombat:latest'
Starting connector 'my-first-connector' with image 'registry.synadia.io/connect-runtime-wombat:latest'
✓ Connector 'my-first-connector' started successfully

[timestamp] INFO Generated message: {"message":"Hello from Wombat"}
[timestamp] INFO Published to events.test
[timestamp] INFO Generated message: {"message":"Hello from Wombat"} 
[timestamp] INFO Published to events.test
...
```

## Step 6: Monitor the Connector

In another terminal, check connector status:

```shell
# List running connectors
connect standalone list

# View logs without following
connect standalone logs my-first-connector

# View detailed container info
docker ps --filter name=my-first-connector
```

## Step 7: Stop the Connector

```shell
# Stop the connector
connect standalone stop my-first-connector

# Remove the container (optional)
connect standalone remove my-first-connector
```

## Customization Ideas

### Change the Message Content

Edit `my-first-connector.connector.yml`:

```yaml
source:
  type: generate
  config:
    mapping: |
      root.id = uuid_v4()
      root.message = "Custom message from " + hostname()
      root.timestamp = now()
      root.counter = counter()
    interval: 2s
```

### Change the Destination

```yaml
producer:
  nats:
    url: "nats://localhost:4222"
  core:
    subject: my.custom.topic
```

### Add Processing

```yaml
steps:
  source:
    # ... source config
  transformer:
    type: mapping
    config:
      mapping: |
        root = this
        root.processed_at = now()
        root.environment = "development"
  producer:
    # ... producer config
```

## Troubleshooting

### Docker Not Running
```
Error: Docker is not available
```
**Solution**: Start Docker Desktop or Docker service

### Port Conflicts
```
Error: address already in use
```
**Solution**: Check if NATS is already running on port 4222

### Container Already Exists
```
Error: container name already exists
```
**Solution**: Remove existing container
```shell
connect standalone remove my-first-connector
```

### Validation Errors
```
Error: failed to parse YAML
```
**Solution**: Check YAML syntax and indentation

## Next Steps

Once you have this working:

1. **Try Different Templates**: `connect standalone template list`
2. **Explore API Integration**: See [API Integration Example](./api-integration.md)
3. **Add Real Data Sources**: See [File Processing Example](./file-processing.md)
4. **Learn About Runtimes**: See [Custom Runtime Example](./custom-runtime.md)

## What You Learned

✅ How to create connectors from templates  
✅ How to validate connector configurations  
✅ How to run connectors in standalone mode  
✅ How to monitor and manage connectors  
✅ How to customize connector behavior

---

**Next Example**: Try [Basic Pipeline](./basic-pipeline.md) to learn data transformation!