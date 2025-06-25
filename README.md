# Synadia Connect

[![Go Report Card](https://goreportcard.com/badge/github.com/synadia-io/connect)](https://goreportcard.com/report/github.com/synadia-io/connect)
[![Coverage Status](https://img.shields.io/badge/coverage-13.2%25-orange.svg)](https://github.com/synadia-io/connect)

Synadia Connect is a NATS-based data pipeline platform that enables seamless data movement between NATS and external systems. Built on top of the NATS ecosystem, Connect provides a simple yet powerful way to create data connectors that can read from various sources and write to different sinks, all while leveraging NATS as the central message broker.

## Deployment Options

Connect supports two deployment modes:

### 🌐 **Managed Service** (Default)
Use Synadia's hosted Connect service for production workloads with full infrastructure management.

### 🐳 **Standalone Mode** (New!)
Run connectors locally using Docker for development, testing, and offline scenarios.

> [!IMPORTANT]
> This project provides the Connect CLI and SDK to interact with the Connect service. The Connect service is currently hosted by Synadia during a private beta. The standalone mode allows you to develop and test connectors locally without requiring access to the hosted service.
> 
> Join us in the #connectors channel on the [NATS Slack](https://slack.nats.io) for questions and discussions!

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Usage](#usage)
- [Documentation](#documentation)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)

## Features

- **Easy Integration**: Simple CLI interface for managing connectors
- **Multiple Runtimes**: Support for Synadia native and Wombat/Benthos runtimes
- **Flexible Connectors**: Create inlets (external → NATS) and outlets (NATS → external)
- **Transformation Support**: Built-in data transformation capabilities
- **Scalable Architecture**: Deploy multiple connector instances for high availability
- **Rich Component Library**: Pre-built components for common data sources and sinks
- **Standalone Mode**: Local development with Docker for offline scenarios

## Installation

### Download Binary

Download the `connect` binary from the [releases page](https://github.com/synadia-io/connect/releases) and place it in your PATH.

### Build from Source

```bash
git clone https://github.com/synadia-io/connect.git
cd connect
task build
task install  # Installs to ~/.local/bin
#Or with basic go build
cd .connect/cmd/connect
go build -o ~/.local/bin/connect.exe
```

### Verify Installation

```bash
# Check version
connect --version

# Ensure correct NATS context (for managed service)
nats context select

# List available components
connect library ls

# Verify Docker for standalone mode
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

### 📋 Available Standalone Commands

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

**Custom     Docker options:**
```shell
# Pass environment variables to connectors
connect standalone run my-connector \
  '--docker-opts=--network host'
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

### 🌐 Managed Service Mode

For production deployments using the hosted service:

```bash
# Create an inlet connector (external source → NATS)
connect connector create my-inlet

# Create an outlet connector (NATS → external sink)
connect connector create my-outlet

# Start with image pulling enabled
connect connector start my-inlet --pull

# Check status
connect connector status my-inlet

# View logs
connect logs
```

For a detailed walkthrough, see the [Getting Started Guide](docs/getting-started.md).

## Architecture

Connect consists of several key components:

### Connector Types

- **Inlets**: Read from external systems and publish to NATS
- **Outlets**: Subscribe to NATS and write to external systems

### Components

Each connector consists of:
- **Source/Consumer**: Where data comes from
- **Sink/Producer**: Where data goes to
- **Transformers** (optional): Data processing steps

### Runtimes

- **Synadia Runtime**: Native NATS-focused components
- **Wombat Runtime**: Benthos-compatible components for broader integrations

## Usage

### CLI Commands

```bash
# Connector management (Managed Service)
connect connector list                    # List all connectors
connect connector create <id>             # Create a new connector
connect connector edit <id>              # Edit existing connector
connect connector delete <id>            # Delete a connector
connect connector start <id> [options]   # Start a connector
connect connector stop <id>              # Stop a connector
connect connector status <id>            # Get connector status

# Standalone mode commands  
connect standalone create <name>          # Create local connector
connect standalone run <name>             # Run locally with Docker
connect standalone logs <name>            # View connector logs

# Library exploration
connect library runtimes                 # List available runtimes
connect library runtime <id>             # Get runtime details
connect library ls [options]             # List components
connect library get <runtime> <kind> <name>  # Get component details

# Monitoring
connect logs                             # View connector logs
```

> 💡 **Tip**: Use `connect standalone` for development and `connect connector` (managed service) for production deployments.

### Configuration

Connectors are configured using YAML specifications. Example:

```yaml
name: my-http-to-nats
description: Reads from HTTP endpoint and publishes to NATS
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
      subject: "events.webhook"
```

## Documentation

### 📚 Getting Started
- [**Getting Started Guide**](docs/getting-started.md) - Create your first connector  
- [**Standalone Mode Guide**](docs/standalone-mode.md) - Complete standalone mode documentation
- [**Examples**](docs/examples/) - Practical examples and tutorials

### 📖 References  
- [**Connector Specification**](spec/schemas/connector-spec.schema.json) - Connector structure details
- [**Connector Reference**](docs/reference/connector.md) - Complete API reference

### 🎯 Quick Links
- [**Quick Start Example**](docs/examples/quick-start.md) - First connector in 5 minutes
- [**API Integration**](docs/examples/api-integration.md) - HTTP to NATS bridge  
- [**Local Development**](docs/examples/local-development.md) - Complete development workflow

## Development

### Prerequisites

- Go 1.22 or later
- Task (task runner)
- NATS Server (for local testing)
- Docker (for standalone mode)

### Building

```bash
# Build all components
task build

# Build specific component
cd connect && task build

# Install to ~/.local/bin
task install
```

### Running Tests

```bash
# Run all tests
task test

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Project Structure

```
connect/
├── cli/          # CLI implementation
├── client/       # Client SDK
├── model/        # Data models (auto-generated)
├── spec/         # Connector specifications
├── builders/     # Fluent API builders
├── convert/      # Conversion utilities
└── docs/         # Documentation
```

## Testing

The project includes comprehensive test coverage:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **Mock Implementations**: Test without external dependencies

Current test coverage: **13.2%** (actively improving)

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./cli -v

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Contributing

We welcome contributions from the community! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code of conduct
- Development workflow
- Submitting pull requests
- Reporting issues

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`task test`)
6. Commit with a descriptive message
7. Push to your branch
8. Open a Pull Request

## Support

- **Slack**: Join #connectors on [NATS Slack](https://slack.nats.io)
- **Issues**: [GitHub Issues](https://github.com/synadia-io/connect/issues)
- **Discussions**: [GitHub Discussions](https://github.com/synadia-io/connect/discussions)

---

Built with ❤️ by [Synadia](https://synadia.com) and the NATS community.