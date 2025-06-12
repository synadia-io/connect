# Synadia Connect

[![License Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/synadia-io/connect)](https://goreportcard.com/report/github.com/synadia-io/connect)
[![Coverage Status](https://img.shields.io/badge/coverage-13.2%25-orange.svg)](https://github.com/synadia-io/connect)

Synadia Connect is a NATS-based data pipeline platform that enables seamless data movement between NATS and external systems. Built on top of the NATS ecosystem, Connect provides a simple yet powerful way to create data connectors that can read from various sources and write to different sinks, all while leveraging NATS as the central message broker.

> [!IMPORTANT]
> This project provides the Connect CLI and SDK to interact with the Connect service. The Connect service is currently hosted by Synadia during a private beta. We're excited to share what we're working on and gather feedback from the community.
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
- [License](#license)

## Features

- **Easy Integration**: Simple CLI interface for managing connectors
- **Multiple Runtimes**: Support for Synadia native and Wombat/Benthos runtimes
- **Flexible Connectors**: Create inlets (external → NATS) and outlets (NATS → external)
- **Transformation Support**: Built-in data transformation capabilities
- **Scalable Architecture**: Deploy multiple connector instances for high availability
- **Rich Component Library**: Pre-built components for common data sources and sinks

## Installation

### Download Binary

Download the `connect` binary from the [releases page](https://github.com/synadia-io/connect/releases) and place it in your PATH.

### Build from Source

```bash
git clone https://github.com/synadia-io/connect.git
cd connect
task build
task install  # Installs to ~/.local/bin
```

### Verify Installation

```bash
# Check version
connect --version

# Ensure correct NATS context
nats context select

# List available components
connect library ls
```

## Quick Start

### 1. Create Your First Connector

```bash
# Create an inlet connector (external source → NATS)
connect connector create my-inlet

# Create an outlet connector (NATS → external sink)
connect connector create my-outlet
```

### 2. Start the Connector

```bash
# Start with image pulling enabled
connect connector start my-inlet --pull

# Check status
connect connector status my-inlet
```

### 3. View Logs

```bash
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
# Connector management
connect connector list                    # List all connectors
connect connector create <id>             # Create a new connector
connect connector edit <id>              # Edit existing connector
connect connector delete <id>            # Delete a connector
connect connector start <id> [options]   # Start a connector
connect connector stop <id>              # Stop a connector
connect connector status <id>            # Get connector status

# Library exploration
connect library runtimes                 # List available runtimes
connect library runtime <id>             # Get runtime details
connect library ls [options]             # List components
connect library get <runtime> <kind> <name>  # Get component details

# Monitoring
connect logs                             # View connector logs
```

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

- [Getting Started Guide](docs/getting-started.md) - Step-by-step tutorial
- [Connector Reference](docs/reference/connector.md) - Detailed connector documentation
- [Component Library](docs/reference/README.md) - Available components and their configurations
- [Connector Specification](spec/schemas/connector-spec.schema.json) - JSON schema for connector definitions
- [API Documentation](docs/api/README.md) - SDK and API reference

## Development

### Prerequisites

- Go 1.22 or later
- Task (task runner)
- NATS Server (for local testing)

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

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

Built with ❤️ by [Synadia](https://synadia.com) and the NATS community.