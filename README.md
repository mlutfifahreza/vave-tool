# Backend - Vave Tool

Go backend service for Vave Tool application.

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

## Architecture Overview

### Layers

1. **cmd/** - Application entry points
   - Each binary has its own subdirectory
   - Minimal logic, mainly wiring dependencies

2. **internal/api/** - HTTP layer
   - **handler/** - HTTP handlers (controllers)
   - **middleware/** - HTTP middleware chain
   - **router/** - Route registration
   - **response/** - Standardized API responses

3. **internal/service/** - Business logic
   - Core application logic
   - Orchestrates repository calls
   - Implements business rules

4. **internal/repository/** - Data access
   - Database operations
   - Query implementations
   - CRUD operations

5. **internal/domain/** - Core domain
   - Domain models/entities
   - Business interfaces
   - Domain errors

6. **internal/pkg/** - Internal utilities
   - Shared internal packages
   - Not exposed outside module

7. **pkg/** - Public packages
   - Can be imported by external projects

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Chi/Gin/Fiber (choose based on preference)
- **Database**: PostgreSQL
- **Migration**: golang-migrate
- **Validation**: validator/v10
- **Authentication**: JWT
- **Logging**: zap/zerolog
- **Testing**: testify

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Environment Setup

Copy the example environment file:

```bash
cp .env.example .env
```

Configure your environment variables:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=vave_tool

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=debug
```

### Installation

Install dependencies:

```bash
go mod download
go mod verify
```

### Database Setup

Run migrations:

```bash
make migrate-up
```

Or manually:

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/vave_tool?sslmode=disable" up
```

### Running the Application

#### Using Make

```bash
make run
```

#### Using Go

```bash
go run cmd/api/main.go
```

#### Using Docker Compose

```bash
docker-compose up
```

## Development

### Running Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests
make test-integration

# With coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Vet code
make vet
```

### Building

```bash
# Build binary
make build

# Build for production
make build-prod

# Build Docker image
make docker-build
```

## API Endpoints

### Health Check

```
GET /health
```

### Authentication

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
```

### Users

```
GET    /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
```

See [API Documentation](docs/api.md) for detailed endpoint specifications.

## Project Conventions

### Naming

- **Files**: lowercase with underscores (e.g., `user_service.go`)
- **Packages**: lowercase, single word (e.g., `handler`, `service`)
- **Interfaces**: noun or adjective (e.g., `UserRepository`, `Validator`)
- **Implementations**: descriptive (e.g., `PostgresUserRepository`)

### Error Handling

- Use custom error types in `domain/errors.go`
- Wrap errors with context using `fmt.Errorf` with `%w`
- Handle errors at appropriate layers

### Testing

- Unit tests alongside source files (`*_test.go`)
- Integration tests in `test/integration/`
- Use table-driven tests
- Aim for 80%+ coverage on business logic

### Security

- All inputs validated using validator
- Passwords hashed with bcrypt
- SQL injection prevention via parameterized queries
- Rate limiting on API endpoints
- CORS configured appropriately
- JWT tokens with expiration

## Dependency Management

Dependencies are managed using Go modules:

```bash
# Add dependency
go get github.com/package/name

# Update dependencies
go get -u ./...

# Tidy dependencies
go mod tidy
```

## Deployment

### Production Build

```bash
make build-prod
```

### Docker Deployment

```bash
docker build -t vave-tool-backend .
docker run -p 8080:8080 --env-file .env vave-tool-backend
```

### Environment Variables

Ensure all required environment variables are set in production:
- Use secret management services
- Never commit secrets to version control

## Monitoring & Logging

- Structured logging with JSON format in production
- Log levels: DEBUG, INFO, WARN, ERROR
- Request ID tracking for traceability
- Health check endpoint for monitoring

## Contributing

1. Follow Go best practices and idioms
2. Write tests for new features
3. Update documentation
4. Run linters before committing
5. Keep commits atomic and well-described

## License

[Add your license here]
