# Vave Tool

Go backend service with REST API and gRPC support.

## 📚 API Documentation

Access the interactive Swagger UI at: **http://localhost:8080/swagger/index.html**

Features:
- Interactive "Try it out" functionality
- Complete endpoint documentation with examples
- Schema definitions and parameter descriptions

**Regenerate docs after changes:**
```bash
make swagger
```

## Quick Start

```bash
# Install dependencies
go mod download

# Create .env file
cp .env.example .env

# Run migrations
make migrate-up

# Make proto
make proto

# Start server
make run

# Access API documentation
# Open http://localhost:8080/swagger/index.html in your browser
```

## 🚀 Traffic Simulator

Generate HTTP traffic for load testing and observability validation.

**Requirements:**
```bash
pip install requests
```

**Basic Usage:**
```bash
# Default: 30 seconds, 5 QPS
python3 script/traffic_simulator.py

# Custom duration and QPS
python3 script/traffic_simulator.py --duration 60 --qps 100

# Short flags
python3 script/traffic_simulator.py -d 120 -q 50

# Custom endpoint
python3 script/traffic_simulator.py -d 30 -q 10 --url http://localhost:8080/api/products/1
```

**Features:**
- ✅ Configurable duration and QPS
- ✅ Real-time progress visualization
- ✅ Concurrent request execution
- ✅ Detailed latency statistics (min, max, mean, median, P95, P99)
- ✅ Status code distribution
- ✅ Success rate calculation
- ✅ Adjustable request timeout

**Example Output:**
```
🚀 Starting traffic simulation...
   Duration: 30 seconds
   QPS:      100 requests/second
   Target:   http://localhost:8080/api/products
   Started:  16:00:19

[██████████████████████████████████████████████████] 3000/3000 | ✓ 2995 | ✗ 5 | Success: 99.8%

📊 Traffic Simulation Complete!
   Total Requests:    3000
   Successful:        2995
   Failed:            5
   Success Rate:      99.83%
   Actual Duration:   30.12s
   Actual QPS:        99.60

📈 Latency Statistics:
   Min:     12.34ms
   Max:     456.78ms
   Mean:    45.67ms
   Median:  42.12ms
   P95:     89.23ms
   P99:     123.45ms
```

## Project Structure

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   ├── worker/
│   │   └── main.go              # Background worker entry point
│   └── migrate/
│       └── main.go              # Database migration tool
│
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTP request handlers
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   └── health.go
│   │   ├── middleware/          # HTTP middlewares
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── ratelimit.go
│   │   ├── router/              # Route definitions
│   │   │   └── router.go
│   │   └── response/            # Response helpers
│   │       └── response.go
│   │
│   ├── service/                 # Business logic layer
│   │   ├── auth.go
│   │   └── user.go
│   │
│   ├── repository/              # Data access layer
│   │   ├── user.go
│   │   └── session.go
│   │
│   ├── domain/                  # Domain models and interfaces
│   │   ├── user.go
│   │   ├── session.go
│   │   └── errors.go
│   │
│   ├── config/                  # Configuration management
│   │   └── config.go
│   │
│   └── pkg/                     # Internal shared packages
│       ├── db/                  # Database connection
│       │   └── postgres.go
│       ├── jwt/                 # JWT utilities
│       │   └── jwt.go
│       ├── hash/                # Password hashing
│       │   └── hash.go
│       ├── validator/           # Input validation
│       │   └── validator.go
│       └── logger/              # Logging utilities
│           └── logger.go
│
├── pkg/                         # Public packages (exportable)
│   └── client/
│       └── client.go
│
├── migrations/                  # Database migrations
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   ├── 000002_add_users.up.sql
│   └── 000002_add_users.down.sql
│
├── scripts/                     # Build and deployment scripts
│   ├── build.sh
│   ├── deploy.sh
│   └── test.sh
│
├── test/                        # Integration and e2e tests
│   ├── integration/
│   └── fixtures/
│
├── docs/                        # API documentation
│   ├── api.md
│   └── swagger.yaml
│
├── .env.example                 # Environment variables template
├── .gitignore
├── go.mod                       # Go module dependencies
├── go.sum
├── Makefile                     # Build automation
├── Dockerfile
├── docker-compose.yml
└── README.md
```