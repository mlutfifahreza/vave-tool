# Quick Start Guide

## 🚀 Start Everything with One Command

```bash
./start.sh
```

This will:
1. Start the observability stack (Prometheus, Loki, Tempo, Grafana)
2. Build the application
3. Start the API server

## 📊 Access the Services

- **API Server**: http://localhost:8080
- **Grafana Dashboard**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **API Metrics**: http://localhost:8080/metrics
- **Swagger Docs**: http://localhost:8080/swagger/

## 🧪 Test the Observability

### 1. Make API Requests

```bash
# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A test product for observability",
    "price": 29.99,
    "stock": 100
  }'

# List products
curl http://localhost:8080/api/products
```

### 2. View in Grafana

1. Open http://localhost:3000
2. Go to **Dashboards** → **Vave Tool API Overview**
3. See metrics, traces, and logs in action!

### 3. Explore Traces

1. In Grafana, click **Explore** (compass icon)
2. Select **Tempo** as data source
3. Click **Search** to see traces
4. Click any trace to see detailed spans

### 4. Explore Logs

1. In Grafana, click **Explore**
2. Select **Loki** as data source  
3. Query: `{service_name="vave-tool-api"}`
4. Click any log entry → Look for **trace_id** → Click to see the full trace!

## 🎯 The Magic: Correlation

When you click a trace in Grafana, you'll see:
- The full request flow (HTTP → Service → Database)
- Timing for each operation
- Related logs with the same **trace_id**

This is the power of the Three Pillars working together!

## 📚 Full Documentation

For detailed information, see [docs/OBSERVABILITY.md](docs/OBSERVABILITY.md)

## 🛑 Cleanup

```bash
make obs-down
```
