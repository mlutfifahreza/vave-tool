# Vave Tool

Go backend service with REST API and gRPC support.

## Quick Start

```bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start server
make run
```

## Development

```bash
make test           # Run tests
make build          # Build binary
make lint           # Run linter
make proto          # Generate proto files
make run            # Start server
```

See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for API details.


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