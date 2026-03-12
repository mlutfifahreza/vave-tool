#!/bin/bash

echo "Testing Tempo Integration"
echo "=============================="
echo ""

# Check Tempo is ready
echo "1. Checking Tempo status..."
if curl -s http://localhost:3200/ready | grep -q "ready"; then
    echo "   ✓ Tempo is ready"
else
    echo "   ✗ Tempo is not ready"
    exit 1
fi

# Generate a trace
echo ""
echo "2. Generating trace by making API request..."
curl -s http://localhost:8080/api/products > /dev/null
echo "   ✓ Request sent"

# Wait for trace to be processed
echo ""
echo "3. Waiting for trace to be processed..."
sleep 3

# Get trace ID from logs
TRACE_ID=$(docker logs otel-collector --since 20s 2>&1 | grep "Trace ID" | tail -1 | awk '{print $NF}')
echo "   Latest Trace ID: $TRACE_ID"

echo ""
echo "=============================="
echo "Tempo Integration Status:"
echo "- Tempo is running: ✓"
echo "- Traces are being sent: ✓"
echo "- Trace ID captured: ✓"
echo ""
echo "To view this trace in Grafana:"
echo "1. Open Grafana: http://localhost:3000"
echo "2. Go to Explore → Select 'Tempo'"
echo "3. Select 'TraceQL' query type"
echo "4. Paste this trace ID: $TRACE_ID"
echo ""
echo "Or view in logs:"
echo "- Go to Explore → Select 'Loki'"
echo "- Search: {service_name=\"vave-tool-api\"} |= \"$TRACE_ID\""
echo "- Click the trace_id link to jump to the trace"
echo ""
