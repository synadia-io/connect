# Local Development Workflow

Complete guide for using standalone mode in your daily development workflow.

## Overview

This guide shows how to set up an efficient development environment using standalone mode for:

✅ **Rapid Prototyping** - Quick connector iteration  
✅ **Local Testing** - Test without external dependencies  
✅ **Team Development** - Consistent development environment  
✅ **CI/CD Integration** - Automated testing pipelines

## Development Environment Setup

### 1. Project Structure

Organize your connector projects:

```
my-connectors/
├── .env                        # Environment variables
├── docker-compose.yml          # Local infrastructure  
├── connectors/
│   ├── data-processor.connector.yml
│   ├── webhook-handler.connector.yml
│   └── file-watcher.connector.yml
├── configs/
│   ├── dev.env
│   ├── test.env
│   └── local.env
├── scripts/
│   ├── start-dev.sh
│   ├── run-tests.sh
│   └── cleanup.sh
└── tests/
    ├── integration/
    └── unit/
```

### 2. Local Infrastructure

Set up supporting services with Docker Compose:

```yaml
# docker-compose.yml
version: '3.8'
services:
  nats:
    image: nats:latest
    ports:
      - "4222:4222"
      - "8222:8222"
    command: ["-js", "-m", "8222"]
    
  mongodb:
    image: mongo:5
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
      
  postgres:
    image: postgres:14
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

Start infrastructure:

```shell
docker-compose up -d
```

### 3. Environment Configuration

Create environment files for different scenarios:

```shell
# configs/dev.env
NATS_URL=nats://localhost:4222
DATABASE_URL=postgres://user:password@localhost:5432/testdb
REDIS_URL=redis://localhost:6379
LOG_LEVEL=debug
ENVIRONMENT=development
```

```shell
# configs/test.env  
NATS_URL=nats://localhost:4222
DATABASE_URL=postgres://user:password@localhost:5432/testdb
REDIS_URL=redis://localhost:6379
LOG_LEVEL=info
ENVIRONMENT=test
```

## Development Workflow

### 1. Create New Connector

```shell
# Create from template
connect standalone create data-processor --template nats-to-mongodb

# Move to project directory
mv data-processor.connector.yml connectors/

# Edit configuration
code connectors/data-processor.connector.yml
```

### 2. Rapid Development Cycle

Create a development script:

```shell
#!/bin/bash
# scripts/dev-cycle.sh

CONNECTOR_NAME=${1:-data-processor}
ENV_FILE=${2:-configs/dev.env}

echo "🔄 Development cycle for: $CONNECTOR_NAME"

# Stop existing connector
echo "Stopping existing connector..."
connect standalone stop $CONNECTOR_NAME 2>/dev/null

# Validate configuration
echo "Validating configuration..."
if ! connect standalone validate $CONNECTOR_NAME --file connectors/$CONNECTOR_NAME.connector.yml; then
    echo "❌ Validation failed"
    exit 1
fi

# Load environment and run
echo "Starting connector with environment: $ENV_FILE"
set -a
source $ENV_FILE
set +a

connect standalone run $CONNECTOR_NAME \
    --file connectors/$CONNECTOR_NAME.connector.yml \
    --follow \
    --remove
```

Use the script:

```shell
# Quick development iteration
./scripts/dev-cycle.sh data-processor configs/dev.env

# Test with different environment
./scripts/dev-cycle.sh data-processor configs/test.env
```

### 3. Live Reload Development

Set up file watching for automatic reload:

```shell
#!/bin/bash
# scripts/watch-and-reload.sh

CONNECTOR_NAME=$1
CONFIG_FILE="connectors/$CONNECTOR_NAME.connector.yml"

echo "👀 Watching $CONFIG_FILE for changes..."

# Install fswatch if needed: brew install fswatch

fswatch -o $CONFIG_FILE | while read f
do
    echo "🔄 File changed, reloading connector..."
    
    # Stop current connector
    connect standalone stop $CONNECTOR_NAME 2>/dev/null
    
    # Validate new configuration
    if connect standalone validate $CONNECTOR_NAME --file $CONFIG_FILE; then
        echo "✅ Configuration valid, restarting..."
        connect standalone run $CONNECTOR_NAME --file $CONFIG_FILE --remove &
    else
        echo "❌ Configuration invalid, not restarting"
    fi
done
```

## Testing Strategy

### 1. Unit Testing Connectors

Create test configurations:

```yaml
# tests/unit/data-processor-test.connector.yml
type: connector
spec:
  id: data-processor-test
  description: Test version with mock data
  runtime_id: wombat
  steps:
    source:
      type: generate
      config:
        mapping: |
          root.id = uuid_v4()
          root.data = "test-data-" + counter()
          root.timestamp = now()
        interval: 1s
        count: 5  # Only generate 5 messages for testing
        
    transformer:
      type: mapping
      config:
        mapping: |
          root = this
          root.test_mode = true
          root.processed_at = now()
          
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: test.output
```

Run tests:

```shell
#!/bin/bash
# scripts/run-tests.sh

echo "🧪 Running connector tests..."

# Start test infrastructure
docker-compose up -d

# Wait for services
sleep 5

# Run unit tests
for test_file in tests/unit/*.connector.yml; do
    test_name=$(basename $test_file .connector.yml)
    echo "Running test: $test_name"
    
    # Run test connector
    connect standalone run $test_name --file $test_file --remove &
    TEST_PID=$!
    
    # Wait for completion
    sleep 10
    
    # Check results
    if connect standalone logs $test_name | grep -q "test-data"; then
        echo "✅ $test_name passed"
    else
        echo "❌ $test_name failed"
    fi
    
    # Cleanup
    connect standalone stop $test_name 2>/dev/null
done

echo "🏁 Tests completed"
```

### 2. Integration Testing

Test complete workflows:

```shell
#!/bin/bash
# scripts/integration-test.sh

echo "🔗 Running integration tests..."

# Start all required connectors
connect standalone run webhook-handler --file connectors/webhook-handler.connector.yml &
sleep 5

connect standalone run data-processor --file connectors/data-processor.connector.yml &
sleep 5

# Send test data
curl -X POST http://localhost:8080/webhook \
    -H "Content-Type: application/json" \
    -d '{"test": true, "data": "integration-test"}'

# Wait for processing
sleep 10

# Verify results
if nats req test.verify '{"check": "integration"}' --timeout 5s; then
    echo "✅ Integration test passed"
else
    echo "❌ Integration test failed"
fi

# Cleanup
connect standalone stop webhook-handler data-processor
```

## Debugging Techniques

### 1. Verbose Logging

Add debug configuration:

```yaml
# Enable debug logging
transformer:
  type: mapping
  config:
    mapping: |
      # Log processing steps
      root = this
      root.debug = {
        "step": "processing",
        "input": this,
        "timestamp": now()
      }
      
      # Your transformation logic here
      root.processed = true
```

### 2. Message Inspection

Create inspection connectors:

```yaml
# debug/message-inspector.connector.yml
type: connector
spec:
  id: message-inspector
  description: Debug message content
  runtime_id: wombat
  steps:
    source:
      type: nats_core
      config:
        urls: ["nats://localhost:4222"]
        subject: ">"  # Subscribe to all messages
        
    transformer:
      type: mapping
      config:
        mapping: |
          root.debug = {
            "subject": meta("nats_subject"),
            "content": this,
            "size_bytes": content().length(),
            "timestamp": now()
          }
          
    producer:
      type: stdout
```

Run inspector:

```shell
connect standalone run message-inspector --file debug/message-inspector.connector.yml --follow
```

### 3. Performance Profiling

Monitor resource usage:

```shell
#!/bin/bash
# scripts/profile-connector.sh

CONNECTOR_NAME=$1

echo "📊 Profiling connector: $CONNECTOR_NAME"

# Start profiling
docker stats $(docker ps --filter name=${CONNECTOR_NAME}_connector --format "{{.Names}}") &
STATS_PID=$!

# Run load test
echo "Generating load..."
for i in {1..1000}; do
    curl -s -X POST http://localhost:8080/webhook \
        -H "Content-Type: application/json" \
        -d '{"test": '$i', "data": "load-test"}' &
done

wait

# Stop profiling
kill $STATS_PID

echo "📈 Profile complete"
```

## Team Collaboration

### 1. Shared Configuration

Use Git for configuration management:

```shell
# .gitignore
.env.local
logs/
*.log
.DS_Store

# docker-compose.override.yml for local customizations
```

### 2. Development Standards

Create team conventions:

```yaml
# .connect-standards.yml
naming:
  connectors: "kebab-case"
  subjects: "dot.notation"
  
structure:
  required_fields:
    - id
    - description  
    - runtime_id
    
  transformer_patterns:
    - Always add timestamp
    - Include processing metadata
    - Use consistent field names
```

### 3. Environment Synchronization

Keep environments in sync:

```shell
#!/bin/bash
# scripts/sync-env.sh

echo "🔄 Synchronizing development environment..."

# Pull latest configurations
git pull origin main

# Update infrastructure
docker-compose pull
docker-compose up -d

# Validate all connectors
for config in connectors/*.connector.yml; do
    connector_name=$(basename $config .connector.yml)
    if ! connect standalone validate $connector_name --file $config; then
        echo "❌ Invalid config: $config"
        exit 1
    fi
done

echo "✅ Environment synchronized"
```

## CI/CD Integration

### 1. GitHub Actions

```yaml
# .github/workflows/test-connectors.yml
name: Test Connectors

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      nats:
        image: nats:latest
        ports:
          - 4222:4222
        options: >-
          --health-cmd "nats server check"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
          
    steps:
    - uses: actions/checkout@v3
    
    - name: Install Connect CLI
      run: |
        curl -L https://github.com/synadia-io/connect/releases/latest/download/connect-linux-amd64 -o connect
        chmod +x connect
        sudo mv connect /usr/local/bin/
        
    - name: Validate Connectors
      run: |
        for config in connectors/*.connector.yml; do
          connector_name=$(basename $config .connector.yml)
          connect standalone validate $connector_name --file $config
        done
        
    - name: Run Integration Tests
      run: |
        ./scripts/run-tests.sh
```

### 2. Pre-commit Hooks

```shell
#!/bin/bash
# .git/hooks/pre-commit

echo "🔍 Pre-commit validation..."

# Validate all connector configurations
for config in connectors/*.connector.yml; do
    if [[ -f $config ]]; then
        connector_name=$(basename $config .connector.yml)
        if ! connect standalone validate $connector_name --file $config; then
            echo "❌ Validation failed for: $config"
            exit 1
        fi
    fi
done

echo "✅ All connectors valid"
```

## Best Practices

### 1. Configuration Management

- **Version Control**: Keep all configs in Git
- **Environment Separation**: Use different files for dev/test/prod
- **Secrets Management**: Use environment variables for sensitive data
- **Documentation**: Document all configuration options

### 2. Development Workflow

- **Feature Branches**: Create branches for new connectors
- **Small Iterations**: Test frequently with small changes
- **Validation First**: Always validate before running
- **Clean Containers**: Remove containers between tests

### 3. Team Coordination

- **Naming Conventions**: Consistent naming across team
- **Shared Infrastructure**: Use docker-compose for common services
- **Code Reviews**: Review connector configurations
- **Documentation**: Keep examples and guides updated

## Troubleshooting

### Common Development Issues

1. **Port Conflicts**: Use docker-compose to manage ports
2. **Configuration Errors**: Use validation early and often
3. **Resource Conflicts**: Clean up containers regularly
4. **Environment Drift**: Use sync scripts to keep environments aligned

### Debug Commands

```shell
# Check running connectors
connect standalone list

# View detailed logs
connect standalone logs <connector> --follow

# Check container status
docker ps --filter name=_connector

# Clean up everything
docker container prune -f
docker image prune -f
```

## Next Steps

1. **Production Deployment**: Transition to managed service
2. **Advanced Monitoring**: Add metrics and alerting
3. **Custom Runtimes**: Create specialized processing engines
4. **Performance Optimization**: Tune for high throughput

---

**What You Learned**:
✅ Complete development environment setup  
✅ Rapid iteration workflows  
✅ Testing strategies  
✅ Debugging techniques  
✅ Team collaboration patterns  
✅ CI/CD integration