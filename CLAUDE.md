# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is the Synadia Connect CLI and SDK - a NATS-based data pipeline tool for connecting data sources and sinks. The project is written in Go and uses a task-based build system.

## Key Commands

### Build and Development
- `task build` - Build the connect binary to `target/connect`
- `task test` - Run all tests with coverage using Ginkgo
- `task install` - Build and install binary to `~/.local/bin/connect`
- `task clean` - Clean up build artifacts in `target/`

### Standalone Mode Commands
Standalone mode allows running connectors locally without NATS services:

#### Connector Management
- `connect standalone validate <file>` - Validate connector definition
- `connect standalone run <file>` - Run connector locally using Docker
- `connect standalone create <name>` - Create new connector definition (creates `<name>.connector.yml`)
- `connect standalone stop <id>` - Stop running connector
- `connect standalone logs <id>` - Show connector logs
- `connect standalone list` - List running connectors

#### Template Management
- `connect standalone template list` - List available templates
- `connect standalone template get <name>` - Generate template file

#### Runtime Management
- `connect standalone runtime list` - List available runtimes
- `connect standalone runtime add <id> <registry>` - Add new runtime configuration
- `connect standalone runtime remove <id>` - Remove runtime configuration
- `connect standalone runtime show <id>` - Show runtime details

Default runtime: `wombat` (registry.synadia.io/connect-runtime-wombat)
Runtime versioning: `wombat` → `registry.synadia.io/connect-runtime-wombat:latest`, `wombat:v1.0.3` → `registry.synadia.io/connect-runtime-wombat:v1.0.3`

### Code Generation
- `task models:generate` - Generate Go models from JSON schemas (requires CONNECT_NODE_LOCATION env var)
- `task models:deps` - Install go-jsonschema dependency for model generation

### Testing
- Tests use Ginkgo/Gomega framework
- Run with: `go run github.com/onsi/ginkgo/v2/ginkgo -r -cover --coverpkg=github.com/synadia-io/connect/... --compilers=4`
- All packages include comprehensive test suites (look for `*_test.go` files)

## Architecture

### Core Components
1. **CLI Layer** (`cli/`): Command-line interface with subcommands for connector, library, and logs operations
2. **Client Layer** (`client/`): HTTP client for communicating with the Connect service
3. **Model Layer** (`model/`): Generated Go structs from JSON schemas for API communication
4. **Spec Layer** (`spec/`): Connector specification definitions and builders
5. **Builders** (`builders/`): Fluent API builders for creating connector configurations
6. **Runtime** (`runtime/`): Core runtime environment for connector execution

### Key Patterns
- **Builder Pattern**: Used extensively in `builders/` for creating connector configurations
- **Generated Code**: Models are auto-generated from JSON schemas - never edit directly
- **Command Pattern**: CLI uses fisk for command structure with separate command files
- **Test Suites**: Each package has its own test suite using Ginkgo

### Important Directories
- `cmd/connect/`: Main entry point and version handling
- `spec/schemas/`: JSON schemas that define the connector specification
- `docs/`: Documentation including getting started guide and reference docs
- `test/`: Shared test utilities and helpers
- `validation/`: Schema validation for standalone mode
- `docker/`: Docker integration for running connectors locally
- `standalone/`: Runtime configuration management for standalone mode

## Development Notes

### Model Generation
- Models in `model/` are generated from JSON schemas - never edit manually
- Use `task models:generate` to regenerate after schema changes
- Requires `CONNECT_NODE_LOCATION` environment variable pointing to schema location

### Testing Requirements
- All new code must have comprehensive test coverage
- Use Ginkgo/Gomega testing framework consistently
- Tests must pass before any changes are committed

### Code Style
- Follow standard Go conventions
- Use builder patterns for complex object construction
- Generated code should not be modified directly
- Prefer dependency injection and interfaces for testability

### Standalone Mode
- Standalone mode provides local connector execution without NATS infrastructure
- Uses Docker for connector runtime isolation with configurable runtime management
- Includes built-in validation and template system
- Supports all connector types (source/sink, consumer/producer, transformers)
- Local ConnectFile format follows same schema as managed service
- Runtime configurations stored in `~/.synadia/connect/standalone/runtimes.json`
- Default runtimes automatically created on first use