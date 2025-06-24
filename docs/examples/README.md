# Standalone Mode Examples

This directory contains practical examples for using Synadia Connect in standalone mode.

## Examples

### 🚀 Available Examples
- [**Quick Start**](./quick-start.md) - Your first connector in 5 minutes
- [**API Integration**](./api-integration.md) - HTTP to NATS bridge
- [**Local Development**](./local-development.md) - Full development setup

### 📋 Coming Soon
Additional examples are in development:
- **Basic Pipeline** - Simple data transformation
- **Template Tour** - Overview of all templates
- **File Processing** - Process files and send to NATS
- **Stream Processing** - Real-time data transformation
- **Database Sync** - NATS to MongoDB pipeline
- **Testing Strategy** - Test connectors effectively
- **Multi-Environment** - Dev/staging/prod configs
- **Custom Runtime** - Add your own runtime
- **High Throughput** - Performance optimization
- **Error Handling** - Robust error management
- **Monitoring** - Observability and logging
- **Security** - Secure configuration practices

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