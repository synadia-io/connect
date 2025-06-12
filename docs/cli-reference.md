# Connect CLI Reference

Complete reference for the Synadia Connect command-line interface.

## Global Options

These options can be used with any command:

```bash
connect [global options] command [command options] [arguments...]
```

| Option | Short | Description | Default | Environment Variable |
|--------|-------|-------------|---------|---------------------|
| `--server URL` | `-s` | NATS server URLs | | `NATS_URL` |
| `--user USER` | | Username or Token | | `NATS_USER` |
| `--password PASSWORD` | | Password | | `NATS_PASSWORD` |
| `--creds FILE` | | User credentials file | | `NATS_CREDS` |
| `--jwt JWT` | | User JWT | | `NATS_JWT` |
| `--seed SEED` | | User seed | | `NATS_SEED` |
| `--context NAME` | | Configuration context | | `NATS_CONTEXT` |
| `--timeout DURATION` | | Response timeout | `5s` | `NATS_TIMEOUT` |
| `--log-level LEVEL` | | Log level (error/warn/info/debug/trace) | `info` | |
| `--help` | `-h` | Show help | | |
| `--version` | | Show version | | |

## Commands

### connector (c)

Manage connectors.

#### connector list (ls)

List all connectors.

```bash
connect connector list
```

Output shows:
- Connector ID
- Description
- Runtime
- Running instances (green)
- Stopped instances (red)

#### connector create/edit

Create or modify a connector.

```bash
connect connector create <id> [options]
connect connector edit <id> [options]
```

Options:
- `--file FILE` (`-f`): Connector specification file

Interactive mode (default):
- Choose from templates
- Edit in preferred editor
- Validate before saving

#### connector get (show, info)

Display connector details.

```bash
connect connector get <id>
```

Shows complete connector configuration in YAML format.

#### connector delete (rm, remove)

Delete a connector.

```bash
connect connector delete <id>
```

> **Note**: Cannot delete running connectors. Stop first.

#### connector start

Start connector instances.

```bash
connect connector start <id> [options]
```

Options:
- `--replicas N`: Number of instances to start (default: 1)
- `--pull`: Pull container image before starting
- `--no-pull`: Skip image pull
- `--env KEY=VALUE`: Set environment variables (can be repeated)
- `--env-file FILE`: Load environment variables from file
- `--placement-tag TAG`: Placement constraints (can be repeated)
- `--timeout DURATION`: Start timeout (default: 2m)

Examples:
```bash
# Start with 3 replicas
connect connector start my-connector --replicas 3

# Start with environment variables
connect connector start my-connector --env API_KEY=secret --env DEBUG=true

# Start with placement constraints
connect connector start my-connector --placement-tag region:us-east --placement-tag env:prod
```

#### connector stop

Stop all instances of a connector.

```bash
connect connector stop <id>
```

#### connector restart

Restart a connector (stop then start).

```bash
connect connector restart <id>
```

#### connector status

Get connector instance status.

```bash
connect connector status <id>
```

Shows:
- Running instances count
- Stopped instances count

#### connector copy (cp)

Copy an existing connector.

```bash
connect connector copy <source-id> <target-id>
```

### library (l)

Explore available components.

#### library runtimes

List available runtimes.

```bash
connect library runtimes
```

Shows:
- Runtime ID
- Name
- Description
- Author

#### library runtime

Get runtime details.

```bash
connect library runtime <id>
```

Shows complete runtime information including:
- Description
- Author details
- Container image
- Metrics configuration

#### library list (ls)

List available components.

```bash
connect library list [options]
```

Options:
- `--runtime ID`: Filter by runtime
- `--kind KIND`: Filter by kind (source/sink/consumer/producer/transformer)
- `--status STATUS`: Filter by status (stable/preview/experimental/deprecated)

Examples:
```bash
# List all components
connect library ls

# List only sources from synadia runtime
connect library ls --runtime synadia --kind source

# List only stable components
connect library ls --status stable
```

#### library get (show)

Get component details.

```bash
connect library get <runtime> <kind> <name>
```

Example:
```bash
connect library get synadia source http_server
```

Shows:
- Component description
- Status
- Configuration fields with:
  - Field names and types
  - Descriptions
  - Default values
  - Examples
  - Constraints

### logs (log)

View connector logs.

```bash
connect logs
```

Streams logs from all running connectors. Press Ctrl+C to stop.

## Configuration Files

### Connector Specification

Connectors are defined using YAML files:

```yaml
name: my-connector
description: Description of what this connector does
runtime: synadia
config:
  source:
    type: component_type
    config:
      field1: value1
      field2: value2
  transformer:  # optional
    type: transformer_type
    config:
      rules: []
  sink:
    type: sink_type
    config:
      target: value
```

### Environment Variables in Configs

Use `${VAR_NAME}` syntax to reference environment variables:

```yaml
config:
  producer:
    type: postgresql
    config:
      connection_string: ${DATABASE_URL}
      password: ${DB_PASSWORD}
```

### Environment Files

Environment files use standard format:

```bash
# .env file
API_KEY=secret-key-value
DEBUG=true
DATABASE_URL=postgresql://user:pass@host:5432/db
```

Load with:
```bash
connect connector start my-connector --env-file .env
```

## Examples

### Create HTTP to NATS Inlet

```bash
# Create interactively
connect connector create webhook-inlet

# Create from file
cat > webhook.yaml << EOF
name: webhook-inlet
description: Receive webhooks and publish to NATS
runtime: synadia
config:
  source:
    type: http_server
    config:
      address: "0.0.0.0:8080"
      path: "/webhook"
  sink:
    type: nats_jetstream
    config:
      subject: "webhooks.received"
      stream: "WEBHOOKS"
EOF

connect connector create webhook-inlet -f webhook.yaml
```

### Create NATS to Database Outlet

```bash
cat > db-writer.yaml << EOF
name: event-writer
description: Write events to PostgreSQL
runtime: synadia
config:
  consumer:
    type: nats_jetstream
    config:
      stream: "EVENTS"
      consumer: "db-writer"
  producer:
    type: postgresql
    config:
      connection_string: \${DATABASE_URL}
      table: "events"
EOF

connect connector create event-writer -f db-writer.yaml
```

### Manage Connector Lifecycle

```bash
# Create
connect connector create my-connector

# Start with 3 instances
connect connector start my-connector --replicas 3 --pull

# Check status
connect connector status my-connector

# View logs
connect logs

# Stop
connect connector stop my-connector

# Delete
connect connector delete my-connector
```

## Tips and Tricks

### 1. Use Contexts

Set up NATS contexts for different environments:

```bash
# Create contexts
nats context add prod --server nats://prod.example.com:4222 --creds prod.creds
nats context add dev --server nats://localhost:4222

# Use specific context
connect --context prod connector list
```

### 2. Debugging

Enable debug logging:

```bash
connect --log-level debug connector start my-connector
```

### 3. Batch Operations

Use shell scripting for batch operations:

```bash
# Start all connectors
for conn in $(connect connector list | grep -oE '^[a-z0-9-]+'); do
  connect connector start $conn
done
```

### 4. Output Formatting

Pipe to standard tools:

```bash
# Get connector count
connect connector list | grep -c "│"

# Extract running connectors
connect connector list | grep "▶" | awk '{print $1}'
```

## Troubleshooting

### Connection Issues

```bash
# Test NATS connection
nats server check

# Use explicit server
connect --server nats://localhost:4222 connector list
```

### Permission Errors

Ensure your NATS credentials have necessary permissions:
- `$JS.API.>` for JetStream operations
- `$SYS.>` for system operations
- Subject permissions for data flow

### Timeout Errors

Increase timeout for slow operations:

```bash
connect --timeout 30s connector start large-connector --replicas 10
```

## See Also

- [Getting Started Guide](getting-started.md)
- [Connector Reference](reference/connector.md)
- [API Documentation](api/README.md)