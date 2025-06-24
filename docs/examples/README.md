# Standalone Mode Examples

This directory contains practical examples for using Synadia Connect in standalone mode.

## Examples

### 🚀 Getting Started
- [**Quick Start**](./quick-start.md) - Your first connector in 5 minutes
- [**Basic Pipeline**](./basic-pipeline.md) - Simple data transformation
- [**Template Tour**](./template-tour.md) - Overview of all templates

### 🔄 Data Pipelines  
- [**File Processing**](./file-processing.md) - Process files and send to NATS
- [**API Integration**](./api-integration.md) - HTTP to NATS bridge
- [**Stream Processing**](./stream-processing.md) - Real-time data transformation
- [**Database Sync**](./database-sync.md) - NATS to MongoDB pipeline

### 🛠️ Development Workflows
- [**Local Development**](./local-development.md) - Full development setup
- [**Testing Strategy**](./testing-strategy.md) - Test connectors effectively  
- [**Multi-Environment**](./multi-environment.md) - Dev/staging/prod configs
- [**Custom Runtime**](./custom-runtime.md) - Add your own runtime

### 🏭 Production Ready
- [**High Throughput**](./high-throughput.md) - Performance optimization
- [**Error Handling**](./error-handling.md) - Robust error management
- [**Monitoring**](./monitoring.md) - Observability and logging
- [**Security**](./security.md) - Secure configuration practices

## Usage

Each example includes:

✅ **Complete Configuration** - Ready-to-run connector files  
✅ **Step-by-Step Instructions** - Clear setup process  
✅ **Expected Output** - What you should see  
✅ **Troubleshooting** - Common issues and solutions  
✅ **Next Steps** - How to extend and customize

## Prerequisites

All examples assume you have:

- Connect CLI with standalone mode
- Docker installed and running  
- Basic familiarity with YAML
- NATS server (for examples that need it)

## Quick Test

Verify your setup works:

```shell
# Test Docker
docker --version

# Test Connect CLI  
connect standalone --help

# Test NATS (if needed)
nats server --jetstream &
```

## Contributing Examples

We welcome community examples! Please:

1. Follow the established format
2. Include complete, working configurations
3. Test thoroughly before submitting
4. Add appropriate documentation

See [Contributing Guide](../../CONTRIBUTING.md) for details.

---

**Get Started**: Try the [Quick Start](./quick-start.md) example first!