# Connect
Connect is a NATS based data pipeline that allows you to easily connect data sources and sinks to NATS. Connect is built on 
top of the NATS ecosystem and uses the NATS server as a message broker. Connect is designed to be easy to use and 
scalable, allowing you to connect multiple data sources and sinks to NATS.

## Deployment Options

Connect supports two deployment modes:

### 🌐 **Managed Service** (Default)
Use Synadia's hosted Connect service for production workloads with full infrastructure management.

### 🐳 **Standalone Mode** (New!)
Run connectors locally using Docker for development, testing, and offline scenarios.

> [!IMPORTANT]
> The managed service is currently hosted by Synadia during a private beta. The standalone mode allows you to 
> develop and test connectors locally without requiring access to the hosted service.
> If you have any questions, feel free to reach out in the #connectors channel on the NATS slack.

## Installation
Download the `connect` binary from the [releases page](https://github.com/synadia-io/connect/releases) and place them in your PATH.

### For Managed Service
Make sure you are using the right nats context:
```shell
nats context select
```

To make sure your binary is correct and everything is up and running, try listing the available components:
```shell
connect library ls
```

### For Standalone Mode
Ensure Docker is installed and running on your system:
```shell
docker --version
```

## Quick Start

### 🚀 Standalone Mode (Recommended for Development)

Get started with standalone mode in 30 seconds:

```shell
# Create a new connector from template
connect standalone create my-connector --template generate

# Validate the connector configuration
connect standalone validate my-connector

# Run the connector locally
connect standalone run my-connector

# View logs
connect standalone logs my-connector --follow

# Stop the connector
connect standalone stop my-connector
```

### 📋 Available Commands

**Connector Management:**
```shell
connect standalone create <name>           # Create connector from template
connect standalone validate <name>         # Validate connector configuration  
connect standalone run <name>              # Run connector in Docker
connect standalone stop <name>             # Stop running connector
connect standalone remove <name>           # Remove connector container
connect standalone logs <name>             # View connector logs
connect standalone list                     # List running connectors
```

**Template Management:**
```shell
connect standalone template list           # List available templates
connect standalone template get <name>     # Get specific template
```

**Runtime Management:**
```shell
connect standalone runtime list            # List available runtimes
connect standalone runtime add <id>        # Add custom runtime
connect standalone runtime remove <id>     # Remove custom runtime  
connect standalone runtime show <id>       # Show runtime details
```

### 🎯 Common Workflows

**Development Workflow:**
```shell
# 1. Create from template
connect standalone create my-data-pipeline --template nats-to-http

# 2. Edit the generated my-data-pipeline.connector.yml file
# 3. Validate your changes
connect standalone validate my-data-pipeline

# 4. Test locally
connect standalone run my-data-pipeline --follow

# 5. Iterate and improve
connect standalone stop my-data-pipeline
# Edit config, then repeat steps 3-5
```

**Testing Different Runtimes:**
```shell
# Add a custom runtime
connect standalone runtime add my-runtime \
  --registry registry.example.com/my-runtime \
  --description "My custom runtime"

# Use in connector spec
# runtime_id: my-runtime:v1.0.0
```

### 🔧 Advanced Usage

**Environment Variables:**
```shell
# Pass environment variables to connectors
connect standalone run my-connector \
  --env DATABASE_URL=postgres://localhost:5432/mydb \
  --env API_KEY=secret123
```

**Custom Configuration:**
```shell
# Override default file locations
connect standalone validate --file ./configs/my-connector.yml
connect standalone run my-connector --file ./configs/my-connector.yml
```

**Container Management:**
```shell
# Remove stopped containers automatically
connect standalone run my-connector --remove

# Follow logs in real-time  
connect standalone logs my-connector --follow

# List all connectors
connect standalone list
```

## Documentation

### 📚 Getting Started
- [**Getting Started Guide**](docs/getting-started.md) - Create your first connector  
- [**Standalone Mode Guide**](docs/standalone-mode.md) - Complete standalone mode documentation
- [**Examples**](docs/examples/) - Practical examples and tutorials

### 📖 References  
- [**Connector Specification**](spec/schemas/connector-spec.schema.json) - Connector structure details
- [**Connector Reference**](docs/reference/connector.md) - Complete API reference
- [**CLI Reference**](docs/cli-reference.md) - Command-line interface guide

### 🎯 Quick Links
- [**Quick Start Example**](docs/examples/quick-start.md) - First connector in 5 minutes
- [**API Integration**](docs/examples/api-integration.md) - HTTP to NATS bridge  
- [**Template Reference**](docs/templates.md) - Available connector templates

## Usage

### Managed Service Mode
You can use the `connect` command to interact with the hosted connect service and manage connectors or explore the library:

```shell
connect --help                    # Show all commands
connect list                      # List connectors  
connect create <name>             # Create connector
connect start <name>              # Start connector
connect stop <name>               # Stop connector
connect library ls                # List available components
```

### Standalone Mode  
Use standalone commands for local development and testing:

```shell
connect standalone --help         # Show standalone commands
connect standalone create <name>  # Create local connector
connect standalone run <name>     # Run locally with Docker
connect standalone logs <name>    # View connector logs
```

> 💡 **Tip**: Use `connect standalone` for development and `connect` (managed service) for production deployments.

## Contributing
We love to get feedback, bug reports, and contributions from our community. If you have any questions or want to
contribute, feel free to reach out in the #connectors channel on the NATS slack.

Take a look at the [Contributing](CONTRIBUTING.md) guide to get started.