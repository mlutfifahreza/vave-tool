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