# Vave Tool - Observability Guide

## Overview

This project implements the **Three Pillars of Observability**:
- **Metrics** - Quantitative measurements (via Prometheus)
- **Logs** - Structured event records (via Loki)
- **Traces** - Request flow tracking (via Tempo)

All data is visualized through **Grafana**, providing a unified dashboard.

## Architecture

```
Application (Go + OpenTelemetry)
    ↓
OpenTelemetry Collector
    ↓
├─→ Prometheus (Metrics)
├─→ Loki (Logs)
└─→ Tempo (Traces)
    ↓
Grafana (Visualization)
```

## Quick Start

### 1. Start the Observability Stack

```bash
make obs-up
```

This starts:
- Prometheus on `localhost:9090`
- Loki on `localhost:3100`
- Tempo on `localhost:3200`
- OpenTelemetry Collector on `localhost:4319` (gRPC) and `localhost:4320` (HTTP)
- Grafana on `localhost:3000`
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`

### 2. Start the Application

```bash
make run
```

The API will be available at `http://localhost:8080`

### 3. Access Grafana

Open `http://localhost:3000` in your browser:
- Username: `admin`
- Password: `admin`

The dashboard "Vave Tool API Overview" should be automatically provisioned.

## Endpoints

### Application Endpoints
- `GET /api/products` - List all products
- `GET /api/products/{id}` - Get product by ID
- `POST /api/products` - Create a product
- `PUT /api/products/{id}` - Update a product
- `DELETE /api/products/{id}` - Delete a product
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics
- `GET /swagger/` - Swagger documentation

### Observability Endpoints
- Grafana: `http://localhost:3000`
- Prometheus: `http://localhost:9090`
- Tempo: `http://localhost:3200`

## The Three Pillars in Action

### 1. Metrics (Prometheus)

The application automatically collects:
- **Request Rate**: Number of requests per second
- **Request Duration**: Response time (p95, p99)
- **Status Codes**: Success/error rates
- **Custom Metrics**: Business-specific counters

View metrics at `http://localhost:8080/metrics` or in Prometheus at `http://localhost:9090`

Example queries:
```promql
# Request rate
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
sum(rate(http_requests_total{status_code=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

### 2. Logs (Loki)

All logs are structured in JSON format with:
- **Timestamp**: ISO8601 format
- **Level**: DEBUG, INFO, WARN, ERROR
- **Message**: Human-readable description
- **Trace ID**: Correlation with traces
- **Context**: Service name, method, path, etc.

Example log entry:
```json
{
  "timestamp": "2026-03-12T10:30:45.123Z",
  "level": "info",
  "service_name": "vave-tool-api",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "msg": "Product created successfully",
  "product_id": "123",
  "duration_seconds": 0.045
}
```

View logs in Grafana by navigating to Explore → Loki

### 3. Traces (Tempo)

Every request is traced through the entire stack:
- HTTP handler
- Service layer
- Repository operations (database calls)
- Cache operations

Each span includes:
- Operation name
- Duration
- Attributes (product_id, cache_hit, etc.)
- Status (OK/Error)

**Viewing traces in Grafana**:

1. **Using Trace IDs from Logs** (Manual copy):
   - Go to Explore → Select **Loki**
   - Query: `{service_name="vave-tool-api"}`
   - Find a log entry and **copy** the `trace_id` value from the log content  
   - Go to Explore → Select **Tempo**
   - Select **"TraceQL"** from the query type dropdown
   - Enter the trace ID and click "Run query"
   
   **Note**: Clickable trace_id links are configured but may encounter "empty ring" errors due to Tempo running in single-instance mode. Direct TraceQL lookup is the recommended method.

2. **Direct Trace Lookup**:
   - Go to Explore → Select **Tempo**
   - **IMPORTANT**: In the query builder dropdown, select **"TraceQL"** (NOT "Search")
   - Paste a trace ID (get from logs or run `./script/test_tempo.sh`)  
   - Click "Run query"
   
   ⚠️ **Common Error**: If you see "error querying live-stores: empty ring", you're either using the wrong query mode or the trace hasn't been flushed to storage yet. Make sure you selected "TraceQL" from the dropdown in the top-left, not "Search". Wait a few seconds for traces to be written to storage.

3. **View Trace Spans**:
   - Each trace shows the complete request flow:
     - `HTTP Request` → Top-level HTTP handler span
     - `ProductService.GetProduct` → Service layer span
     - `Repository.GetProductByID` → Database layer span
     - Redis cache operations spans

## Correlation IDs - The Magic Connection

The **trace_id** is automatically injected into all logs, enabling you to:

1. **Click a trace in Grafana** → See all related logs
2. **Click a log entry** → Jump to the full trace
3. **See the complete request journey** from HTTP request to database query

This is configured in Grafana's datasources:
- Tempo → logs uses `trace_id` to find related logs
- Loki → traces uses regex `trace_id=(\w+)` to extract and link to traces

## Common Workflows

### Debugging a Slow Request

1. Go to Grafana → Dashboard → "Vave Tool API Overview"
2. Look at "Request Duration" panel for spikes
3. Click on the time range with high latency
4. Navigate to Explore → Tempo
5. Find slow traces → Click to see span details
6. Identify which operation is slow (cache, database, etc.)
7. Click "Logs for this span" to see related log entries

### Investigating an Error

1. Check the "Success Rate" gauge in Grafana dashboard
2. Navigate to Explore → Loki
3. Filter by level: `{service_name="vave-tool-api"} |= "level=\"error\""`
4. Find the error log → Note the `trace_id`
5. Click "Tempo" link in the log to see the full trace
6. Analyze the trace to understand the error context

### Monitoring Application Health

1. Open Grafana Dashboard → "Vave Tool API Overview"
2. Monitor:
   - **Request Rate**: Should be stable
   - **Request Duration**: p95 < 200ms, p99 < 500ms
   - **Success Rate**: Should be > 99%
   - **Recent Logs**: Check for errors/warnings

## Making Test Requests

Generate some traffic to see observability in action:

```bash
# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "price": 29.99,
    "stock": 100
  }'

# List products
curl http://localhost:8080/api/products

# Get a product (replace {id} with actual ID)
curl http://localhost:8080/api/products/{id}

# Update a product
curl -X PUT http://localhost:8080/api/products/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "description": "Updated description",
    "price": 39.99,
    "stock": 50
  }'

# Delete a product
curl -X DELETE http://localhost:8080/api/products/{id}
```

After making requests:
1. Check metrics at `http://localhost:8080/metrics`
2. View traces in Grafana → Explore → Tempo
3. View logs in Grafana → Explore → Loki
4. See the dashboard at Grafana → Dashboards → "Vave Tool API Overview"

## Makefile Commands

```bash
make obs-up      # Start observability stack
make obs-down    # Stop observability stack
make obs-logs    # View container logs
make run         # Run the application
make build       # Build the application
make test        # Run tests
```

## Cleanup

```bash
make obs-down
```

To remove all data:
```bash
docker-compose down -v
```

## Customization

### Adding Custom Metrics

```go
meter := otel.Meter("vave-tool-api")
counter, _ := meter.Int64Counter("custom_operation_total")
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("operation", "custom"),
))
```

### Adding Custom Spans

```go
ctx, span := observability.StartSpan(ctx, "CustomOperation",
    attribute.String("key", "value"),
)
defer span.End()

// Your code here

if err != nil {
    observability.RecordError(span, err, "Operation failed")
}
```

### Adding Structured Logs

```go
logger.Info(ctx, "Custom operation completed",
    zap.String("key", "value"),
    zap.Int("count", 42),
)
```

## Troubleshooting

### Application can't connect to OpenTelemetry Collector

Check if the collector is running:
```bash
docker ps | grep otel-collector
```

View collector logs:
```bash
docker logs otel-collector
```

### No data in Grafana

1. Verify data sources in Grafana Settings → Data Sources
2. Check if applications are sending data:
   ```bash
   # Check Prometheus targets
   curl http://localhost:9090/api/v1/targets
   
   # Check Tempo readiness
   curl http://localhost:3200/ready
   ```

### Traces not showing up

Ensure the application is configured to send traces to the correct endpoint:
- gRPC: `localhost:4319`
- HTTP: `localhost:4320`

**To verify traces are working:**
```bash
# Run the test script
./script/test_tempo.sh
```

### "Empty ring" error in Tempo

**Error message**: `error querying live-stores in Querier.FindTraceByID: error finding partition ring replicas: empty ring`

**Cause**: This error occurs when Tempo tries to query ingesters (live data) in single-instance mode. The current setup runs Tempo in local/dev mode without a distributed ring infrastructure.

**Solutions**:

1. **Wait for trace flush**: Traces may not be immediately queryable. Wait 5-10 seconds after generating a request, then query again.

2. **Use TraceQL for direct lookupup** (Recommended):
   - In Grafana → Explore → Tempo, select **"TraceQL"** from the query type dropdown (NOT "Search")
   - Get a trace ID from:
     - Application logs (look for `trace_id` field)
     - Loki logs (copy the trace_id value from log content)
     - Run `./script/test_tempo.sh`
   - Enter the trace ID in the query field
   - Click "Run query"

3. **Clickable links from Loki**: While configured, trace_id links from Loki logs may encounter this error. Use manual copy-paste method instead.

**Why this happens**: In single-instance mode, Tempo doesn't maintain a memberlist ring for coordinating distributed queries. Direct trace ID lookups via TraceQL work reliably, but some query paths (like the API endpoint used by derived field links) attempt ingester queries that require ring coordination.

**Note**: This doesn't affect trace storage or viewing by ID. All traces are stored and fully functional when accessed directly via TraceQL.

## Learn More

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Prometheus Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [LogQL (Loki Query Language)](https://grafana.com/docs/loki/latest/logql/)
- [Tempo Documentation](https://grafana.com/docs/tempo/latest/)
- [Grafana Documentation](https://grafana.com/docs/)
