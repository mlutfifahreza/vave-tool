#!/bin/bash

set -e

echo "🚀 Starting Vave Tool with Observability Stack..."
echo ""

echo "📊 Step 1: Starting observability services (Prometheus, Loki, Tempo, Grafana)..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to be ready (30 seconds)..."
sleep 30

echo ""
echo "✅ Observability stack is running!"
echo ""
echo "🔗 Access Points:"
echo "   - Grafana:    http://localhost:3000 (admin/admin)"
echo "   - Prometheus: http://localhost:9090"
echo "   - Tempo:      http://localhost:3200"
echo "   - Loki:       http://localhost:3100"
echo ""

echo "🏗️  Step 2: Building the application..."
go build -o bin/api cmd/api/main.go

echo ""
echo "🎯 Step 3: Starting the API server..."
echo ""
echo "📝 Access the API at: http://localhost:8080"
echo "📊 Metrics endpoint:  http://localhost:8080/metrics"
echo "📚 API Documentation: http://localhost:8080/swagger/"
echo ""
echo "Press Ctrl+C to stop the application"
echo ""

./bin/api
