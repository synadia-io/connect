# Standalone Mode Guide

Standalone mode allows you to run Synadia Connect connectors locally using Docker, perfect for development, testing, and offline scenarios.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Command Reference](#command-reference)
- [Templates](#templates)
- [Runtime Management](#runtime-management)
- [Configuration](#configuration)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Overview

Standalone mode provides a complete local development environment for Synadia Connect:

✅ **No Network Required** - Work offline without internet connectivity  
✅ **Docker-Based** - Leverages Docker for consistent runtime environments  
✅ **Template System** - Quick start with pre-built connector templates  
✅ **Multiple Runtimes** - Support for different connector runtime engines  
✅ **Development Workflow** - Rapid iteration and testing cycle

## Prerequisites

### System Requirements

- **Docker**: Version 20.0+ installed and running
- **Connect CLI**: Latest version with standalone mode support
- **Disk Space**: ~500MB for runtime images and containers

### Verify Installation

```shell
# Check Docker
docker --version
docker ps

# Check Connect CLI
connect --version
connect standalone --help
```

## Quick Start

### 1. Create Your First Connector

```shell
# Create a connector from template
connect standalone create my-first-connector --template generate-to-nats

# This creates: my-first-connector.connector.yml
```

### 2. Validate Configuration

```shell
# Validate the connector configuration
connect standalone validate my-first-connector

# Output: ✓ Connector 'my-first-connector' is valid
```

### 3. Run Locally

```shell
# Run the connector in Docker
connect standalone run my-first-connector

# Output: 
# Using runtime 'wombat' resolved to image 'registry.synadia.io/connect-runtime-wombat:latest'
# Starting connector 'my-first-connector' with image '...'
# ✓ Connector 'my-first-connector' started successfully
```

### 4. Monitor and Manage

```shell
# View real-time logs
connect standalone logs my-first-connector --follow

# List running connectors
connect standalone list

# Stop when done
connect standalone stop my-first-connector
```

## Core Concepts

### Connector Files

Connectors are defined in YAML files with the `.connector.yml` extension:

```yaml
# my-connector.connector.yml
type: connector
spec:
  id: my-connector
  description: My first connector
  runtime_id: wombat
  steps:
    source:
      type: generate
      config:
        mapping: 'root.message = "Hello World"'
        interval: 5s
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: events.test
```

### Naming Convention

- **Connector Name**: `my-connector`
- **File Name**: `my-connector.connector.yml`  
- **Container Name**: `my-connector_connector`

### Runtime Engines

Standalone mode supports multiple runtime engines:

- **wombat** (default): High-performance streaming processor
- **Custom runtimes**: Add your own runtime engines

## Command Reference

### Connector Management

#### `create` - Create New Connector
```shell
connect standalone create <name> [options]

Options:
  --template <name>    Template to use (default: first available)
  --file <path>        Override output file path

Examples:
  connect standalone create my-app --template nats-to-http
  connect standalone create data-sync --file ./configs/sync.yml
```

#### `validate` - Validate Configuration
```shell
connect standalone validate <name> [options]

Options:
  --file <path>        Override input file path

Examples:
  connect standalone validate my-app
  connect standalone validate --file ./configs/custom.yml
```

#### `run` - Execute Connector
```shell
connect standalone run <name> [options]

Options:
  --follow            Follow logs after starting
  --remove            Remove container when it stops
  --env KEY=VALUE     Set environment variables
  --docker-opts <docker options>     Set environment variables
  --image <image>     Override runtime image

Examples:
  connect standalone run my-app --follow
  connect standalone run my-app --env DEBUG=true --env PORT=8080
  connect standalone run my-app --docker-opts '--network=host'
  connect standalone run my-app --image custom-runtime:v1.0.0
   
```

#### `stop` - Stop Running Connector
```shell
connect standalone stop <name>

Examples:
  connect standalone stop my-app
```

#### `remove` - Remove Container
```shell
connect standalone remove <name>

Examples:
  connect standalone remove my-app
```

#### `logs` - View Logs
```shell
connect standalone logs <name> [options]

Options:
  --follow, -f        Follow logs in real-time

Examples:
  connect standalone logs my-app
  connect standalone logs my-app --follow
```

#### `list` - List Connectors
```shell
connect standalone list

# Shows running connectors with status
```

### Template Management

#### `template list` - List Templates
```shell
connect standalone template list

# Shows all available templates
```

#### `template get` - Extract Template
```shell
connect standalone template get <name> [options]

Options:
  --output <path>     Output file path

Examples:
  connect standalone template get nats-to-http
  connect standalone template get generate --output ./my-template.yml
```

### Runtime Management

#### `runtime list` - List Runtimes
```shell
connect standalone runtime list

# Shows available runtime engines
```

#### `runtime add` - Add Custom Runtime
```shell
connect standalone runtime add <id> [options]

Options:
  --registry <url>        Container registry URL
  --description <text>    Runtime description
  --author <name>         Runtime author

Examples:
  connect standalone runtime add my-runtime \
    --registry registry.example.com/my-runtime \
    --description "Custom processing runtime"
```

#### `runtime remove` - Remove Runtime
```shell
connect standalone runtime remove <id>

Examples:
  connect standalone runtime remove my-runtime
```

#### `runtime show` - Show Runtime Details
```shell
connect standalone runtime show <id>

Examples:
  connect standalone runtime show wombat
```

## Templates

### Available Templates

| Template | Description | Use Case |
|----------|-------------|----------|
| `generate-to-nats` | Generate test data to NATS | Development, testing |
| `nats-to-http` | Stream NATS messages to HTTP | API integration |
| `nats-to-stream` | Bridge NATS subjects/streams | Message routing |
| `nats-to-mongodb` | Store NATS messages in MongoDB | Data persistence |

### Using Templates

```shell
# List available templates
connect standalone template list

# Create from specific template
connect standalone create web-forward--template nats-to-http

# Extract template for customization
connect standalone template get nats-to-http --output base-template.yml
```

### Template Structure

Templates provide pre-configured connector definitions:

```yaml
spec:
    description: 'Connector: nats-to-http'
    runtime_id: wombat
    steps:
        consumer:
            core:
                queue: workers
                subject: events.http
            nats:
                url: nats://localhost:4222
        sink:
            config:
                url: http://localhost:3000/webhook
                verb: POST
            type: http_client
type: connector
```

## Runtime Management

### Default Runtime

The `wombat` runtime is included by default:

```shell
connect standalone runtime show wombat

# Output:
# ╭──────────────────────────────────────────╮
# │ Wombat Runtime                           │
# ├─────────────────┬────────────────────────┤
# │ Id              │ wombat                 │
# │ Registry        │ registry.synadia.io... │
# │ Description     │ Default Wombat-based...│
# ╰─────────────────┴────────────────────────╯
```

### Adding Custom Runtimes

```shell
# Add a custom runtime
connect standalone runtime add benthos \
  --registry benthos/benthos \
  --description "Benthos streaming processor"

# Use in connector spec
# runtime_id: benthos:v4.0.0
```

### Runtime Versioning

Specify runtime versions in connector specs:

```yaml
spec:
  runtime_id: wombat:v1.2.3  # Use specific version
  runtime_id: wombat         # Use latest version
```

## Configuration

### Environment Variables

Pass environment variables to connectors:

```shell
# Single variable
connect standalone run my-app --env DATABASE_URL=postgres://localhost/db

# Multiple variables  
connect standalone run my-app \
  --env DATABASE_URL=postgres://localhost/db \
  --env API_KEY=secret123 \
  --env DEBUG=true
```

### Custom Docker options

Pass additional docker options for local testing

```shell
connect standalone run my-app --docker-opts='--network host -p 8080:8080'
```

### File Locations

By default, standalone mode looks for files following this pattern:

- **Connector file**: `<name>.connector.yml`
- **Runtime config**: `~/.synadia/connect/standalone/`

Override with flags:

```shell
# Custom file location
connect standalone validate --file ./configs/my-connector.yml
connect standalone run my-connector --file ./configs/my-connector.yml
```

### Container Management

Control container lifecycle:

```shell
# Auto-remove container when stopped
connect standalone run my-app --remove

# Keep container for debugging
connect standalone run my-app

# Manual cleanup
connect standalone remove my-app
```

## Examples

### Example 1: Data Generation Pipeline

Create a connector that generates test data:

```shell
# 1. Create from template
connect standalone create data-generator --template generate

# 2. Customize the generated file
cat > data-generator.connector.yml << 'EOF'
type: connector
spec:
  id: data-generator
  description: Generate sample user events
  runtime_id: wombat
  steps:
    source:
      type: generate
      config:
        mapping: |
          root.user_id = uuid_v4()
          root.event = "user.signup"
          root.timestamp = now()
          root.email = "user" + random_int(1,1000).string() + "@example.com"
        interval: 2s
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: events.users
EOF

# 3. Run the connector
connect standalone run data-generator --follow
```

### Example 2: HTTP to NATS Bridge

Create a webhook receiver:

```shell
# 1. Create HTTP to NATS bridge
connect standalone create webhook-bridge --template http-to-nats

# 2. Edit configuration
cat > webhook-bridge.connector.yml << 'EOF'
type: connector
spec:
  id: webhook-bridge
  description: Receive webhooks and forward to NATS
  runtime_id: wombat
  steps:
    source:
      type: http_server
      config:
        address: "0.0.0.0:8080"
        path: /webhook
    transformer:
      type: mapping
      config:
        mapping: |
          root.received_at = now()
          root.source = "webhook"
          root.payload = this
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: webhooks.received
EOF

# 3. Run the bridge
connect standalone run webhook-bridge

# 4. Test with curl
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"event": "test", "data": {"key": "value"}}'
```

### Example 3: Multi-Runtime Setup

Use different runtimes for different connectors:

```shell
# 1. Add custom runtime
connect standalone runtime add benthos \
  --registry benthos/benthos \
  --description "Benthos processor"

# 2. Create connector using custom runtime  
cat > custom-processor.connector.yml << 'EOF'
type: connector
spec:
  id: custom-processor
  description: Use custom Benthos runtime
  runtime_id: benthos:latest
  steps:
    source:
      type: nats_core
      config:
        urls: ["nats://localhost:4222"]
        subject: input.data
    transformer:
      type: mapping
      config:
        mapping: |
          root = this
          root.processed_by = "benthos"
          root.processed_at = now()
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: output.processed
EOF

# 3. Run with custom runtime
connect standalone run custom-processor
```

## Troubleshooting

### Common Issues

#### Docker Not Running
```
Error: failed to run connector: Docker is not available
```
**Solution**: Start Docker service
```shell
# macOS/Windows: Start Docker Desktop
# Linux: 
sudo systemctl start docker
```

#### Port Already in Use
```
Error: bind: address already in use
```
**Solution**: Change port in connector config or stop conflicting service
```yaml
source:
  type: http_server
  config:
    address: "0.0.0.0:8081"  # Use different port
```

#### Container Name Conflicts
```
Error: container name already exists
```
**Solution**: Remove existing container
```shell
connect standalone remove my-connector
connect standalone run my-connector
```

#### Runtime Image Not Found
```
Error: failed to resolve runtime: runtime 'custom' not found
```
**Solution**: Add the runtime first
```shell
connect standalone runtime add custom --registry registry.example.com/custom
```

### Debug Mode

Enable verbose logging:

```shell
# View detailed logs
connect standalone logs my-connector --follow

# Check container status
docker ps -a --filter name=my-connector

# Inspect container
docker inspect my-connector_connector
```

### Validation Issues

Common validation errors:

```shell
# Check syntax
connect standalone validate my-connector

# Common fixes:
# - Ensure proper YAML indentation
# - Check required fields (id, runtime_id, steps)  
# - Verify step types are valid
```

## Best Practices

### Development Workflow

1. **Start with Templates**: Use built-in templates as starting points
2. **Validate Early**: Run validation after each change
3. **Use Follow Mode**: Monitor logs with `--follow` during development
4. **Clean Containers**: Use `--remove` to avoid container buildup
5. **Version Control**: Keep connector files in version control

### Configuration Management

```shell
# Organize by environment
./connectors/
  ├── dev/
  │   ├── data-generator.connector.yml
  │   └── webhook-bridge.connector.yml
  ├── staging/
  │   └── production-sync.connector.yml
  └── templates/
      └── base-template.yml

# Use environment-specific configs
connect standalone run data-generator --file ./connectors/dev/data-generator.connector.yml
```

### Performance Tips

1. **Resource Limits**: Monitor Docker resource usage
2. **Log Rotation**: Use log rotation for long-running connectors
3. **Image Cleanup**: Periodically clean unused Docker images
4. **Batch Operations**: Use appropriate batch sizes for high-throughput scenarios

### Security Considerations

1. **Environment Variables**: Use env vars for sensitive configuration
2. **Network Access**: Be mindful of exposed ports
3. **Image Sources**: Use trusted container registries
4. **File Permissions**: Secure connector configuration files

```shell
# Secure approach
connect standalone run my-connector \
  --env DATABASE_PASSWORD="$(cat /secure/db-password)" \
  --env API_TOKEN="$(vault kv get -field=token secret/api)"
```

### Production Transition

When moving from standalone to managed service:

1. **Test Compatibility**: Ensure connector works in both modes
2. **Environment Parity**: Match runtime versions
3. **Configuration Migration**: Adapt environment-specific settings
4. **Monitoring Setup**: Plan observability strategy

```shell
# Development (standalone)
connect standalone run my-connector --env ENV=dev

# Production (managed service) 
connect connector create my-connector --env ENV=prod
connect connector start my-connector
```

## Getting Help

- **Command Help**: `connect standalone --help`
- **Community**: Join #connectors on NATS Slack  
- **Documentation**: See [Connect Reference](./reference/connector.md)
- **Issues**: Report bugs on GitHub

---

**Next Steps**: 
- Try the [Getting Started Guide](./getting-started.md)
- Explore [Connector Examples](./examples/)
- Read [Connector Reference](./reference/connector.md)