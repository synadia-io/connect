# API Integration Example

Build an HTTP-to-NATS bridge that receives webhooks and forwards them to NATS streams.

## Overview

This example demonstrates:
- Setting up an HTTP server to receive webhooks
- Transforming incoming HTTP requests
- Publishing processed data to NATS
- Error handling and monitoring

## Use Cases

✅ **Webhook Processing** - Receive GitHub, Stripe, or custom webhooks  
✅ **API Gateway** - Bridge HTTP APIs to event streams  
✅ **Microservice Integration** - Connect REST services to event-driven architecture  
✅ **Real-time Notifications** - Transform API calls into events

## Step 1: Create the Webhook Bridge

```shell
# Create connector from HTTP template
connect standalone create webhook-bridge --template http-to-nats
```

## Step 2: Configure the Bridge

Edit `webhook-bridge.connector.yml`:

```yaml
# Synadia Connect Connector Definition
type: connector
spec:
  id: webhook-bridge
  description: HTTP to NATS webhook bridge
  runtime_id: wombat
  steps:
    source:
      type: http_server
      config:
        address: "0.0.0.0:8080"
        path: "/webhook"
        timeout: "30s"
        
    transformer:
      type: mapping
      config:
        mapping: |
          # Add metadata to incoming webhook
          root.webhook_id = uuid_v4()
          root.received_at = now()
          root.source_ip = meta("http_request_remote_addr")
          root.user_agent = meta("http_request_user_agent")
          
          # Extract headers we care about
          root.headers.content_type = meta("http_request_header_content_type")
          root.headers.signature = meta("http_request_header_x_signature")
          
          # Include original payload
          root.payload = this
          
          # Add routing information
          root.event_type = this.type | "unknown"
          root.routing_key = "webhooks." + root.event_type
          
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: webhooks.incoming
```

## Step 3: Validate and Run

```shell
# Validate configuration
connect standalone validate webhook-bridge

# Start the bridge
connect standalone run webhook-bridge --follow
```

Expected output:
```
Using runtime 'wombat' resolved to image 'registry.synadia.io/connect-runtime-wombat:latest'
Starting connector 'webhook-bridge'...
✓ Connector 'webhook-bridge' started successfully

[timestamp] INFO HTTP server listening on :8080
[timestamp] INFO Waiting for webhook requests on /webhook
```

## Step 4: Test the Webhook

### Send Test Webhook

```shell
# Simple test webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Signature: test-signature" \
  -d '{
    "type": "user.created",
    "user_id": "12345",
    "email": "user@example.com",
    "timestamp": "2024-01-15T10:30:00Z"
  }'
```

### GitHub-style Webhook

```shell
# Simulate GitHub webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-GitHub-Delivery: 12345-67890" \
  -d '{
    "type": "push",
    "repository": {
      "name": "my-repo",
      "full_name": "user/my-repo"
    },
    "commits": [
      {
        "id": "abc123",
        "message": "Update README",
        "author": {
          "name": "Developer",
          "email": "dev@example.com"
        }
      }
    ]
  }'
```

### Stripe-style Webhook

```shell
# Simulate Stripe webhook  
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "Stripe-Signature: t=1234567890,v1=signature_here" \
  -d '{
    "type": "payment_intent.succeeded",
    "data": {
      "object": {
        "id": "pi_1234567890",
        "amount": 2000,
        "currency": "usd",
        "status": "succeeded"
      }
    }
  }'
```

## Step 5: Monitor Incoming Data

### View Connector Logs

```shell
# Follow webhook processing logs
connect standalone logs webhook-bridge --follow
```

### Subscribe to NATS Messages

```shell
# Subscribe to see processed webhooks
nats sub "webhooks.>"

# Subscribe to specific event types
nats sub "webhooks.user.created"
nats sub "webhooks.push"
nats sub "webhooks.payment_intent.succeeded"
```

## Advanced Configuration

### Multiple Endpoints

Handle different webhook types on different paths:

```yaml
# webhook-multi.connector.yml
spec:
  steps:
    source:
      type: http_server
      config:
        address: "0.0.0.0:8080"
        
    transformer:
      type: mapping
      config:
        mapping: |
          # Route based on URL path
          root.webhook_id = uuid_v4()
          root.received_at = now()
          root.endpoint = meta("http_request_url")
          
          # Extract service from path
          let path = meta("http_request_url").trim_prefix("/webhook/")
          root.service = path.split("/").index(0) | "unknown"
          
          root.payload = this
          
          # Dynamic subject routing
          root.subject = "webhooks." + root.service + "." + (this.type | "event")
          
    producer:
      nats:
        url: "nats://localhost:4222"
      core:
        subject: ${! json("subject") }
```

Test different endpoints:

```shell
# GitHub webhooks
curl -X POST http://localhost:8080/webhook/github -d '{"type":"push"}'

# Stripe webhooks  
curl -X POST http://localhost:8080/webhook/stripe -d '{"type":"payment_intent.succeeded"}'

# Custom service webhooks
curl -X POST http://localhost:8080/webhook/myservice -d '{"type":"user.updated"}'
```

### Authentication & Security

Add webhook signature validation:

```yaml
transformer:
  type: mapping
  config:
    mapping: |
      # Validate webhook signature (example for GitHub)
      let signature = meta("http_request_header_x_hub_signature_256")
      let payload = content().string()
      let secret = env("WEBHOOK_SECRET")
      let expected = "sha256=" + payload.hash_hmac_sha256(secret).encode_hex()
      
      # Fail if signature doesn't match
      if signature != expected {
        throw("Invalid webhook signature")
      }
      
      root.webhook_id = uuid_v4()
      root.verified = true
      root.payload = this
```

Run with secret:

```shell
connect standalone run webhook-bridge \
  --env WEBHOOK_SECRET=your-secret-key-here
```

### Error Handling

Add robust error handling:

```yaml
transformer:
  type: mapping
  config:
    mapping: |
      # Wrap processing in try-catch
      root = {}
      
      try {
        root.webhook_id = uuid_v4()
        root.received_at = now()
        root.payload = this
        root.status = "processed"
        
        # Validate required fields
        if !this.type.exists() {
          throw("Missing required field: type")
        }
        
        root.event_type = this.type
        
      } catch e {
        # Handle errors gracefully
        root.error = e
        root.status = "error"
        root.raw_payload = content().string()
        root.event_type = "webhook.error"
      }
```

## Production Considerations

### SSL/TLS Termination

For production, use a reverse proxy:

```yaml
# docker-compose.yml
version: '3.8'
services:
  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./certs:/etc/nginx/certs
    depends_on:
      - webhook-bridge
      
  webhook-bridge:
    # Your connector container
    ports:
      - "8080:8080"
```

### Rate Limiting

Add rate limiting to prevent abuse:

```yaml
source:
  type: http_server
  config:
    address: "0.0.0.0:8080"
    rate_limit: "100req/s"
    timeout: "30s"
```

### Monitoring & Metrics

Track webhook processing:

```yaml
transformer:
  type: mapping
  config:
    mapping: |
      root = this
      root.webhook_id = uuid_v4()
      root.received_at = now()
      
      # Add metrics labels
      root.metrics = {
        "webhook_type": this.type | "unknown",
        "source_service": meta("http_request_header_user_agent").re_find_all("([A-Za-z]+)").index(0) | "unknown",
        "response_time_ms": (now() - root.received_at).format_timestamp_unix_milli()
      }
```

## Troubleshooting

### Port Already in Use

```
Error: bind: address already in use
```

**Solution**: Change port or stop conflicting service

```yaml
source:
  type: http_server
  config:
    address: "0.0.0.0:8081"  # Use different port
```

### Webhook Not Received

1. **Check connector logs**: `connect standalone logs webhook-bridge`
2. **Verify endpoint**: Test with `curl` first
3. **Check firewall**: Ensure port 8080 is accessible
4. **Validate JSON**: Ensure webhook payload is valid JSON

### Processing Errors

```
Error: failed to process webhook
```

**Debug steps**:

```shell
# Check detailed logs
connect standalone logs webhook-bridge --follow

# Test with minimal payload
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{}'

# Validate transformer mapping
connect standalone validate webhook-bridge
```

## Next Steps

1. **Add Authentication**: Implement webhook signature validation
2. **Scale Horizontally**: Run multiple instances behind load balancer
3. **Add Persistence**: Store failed webhooks for retry
4. **Custom Routing**: Route to different NATS subjects based on content
5. **Monitoring**: Add metrics and alerting

## Related Examples

- [Stream Processing](./stream-processing.md) - Process webhook data further
- [Database Sync](./database-sync.md) - Store webhook data in databases
- [Error Handling](./error-handling.md) - Robust error management
- [Monitoring](./monitoring.md) - Production monitoring setup

---

**What You Learned**:
✅ HTTP server configuration  
✅ Webhook processing patterns  
✅ Data transformation with mapping  
✅ Dynamic routing  
✅ Error handling strategies