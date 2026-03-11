# Vave Tool - Backend API

A Go-based backend service that provides both REST and gRPC APIs for managing products.

## Features

- ✅ RESTful API endpoints
- ✅ gRPC service
- ✅ PostgreSQL database with migrations
- ✅ Clean architecture (Handler → Service → Repository → Domain)
- ✅ CORS support

## API Endpoints

### REST API

**Base URL**: `http://localhost:8080`

#### List Products
```http
GET /api/products
```

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Product Name",
      "description": "Description",
      "price": 99.99,
      "stock_quantity": 100,
      "category": "Electronics",
      "sku": "PROD-001",
      "is_active": true,
      "created_at": "2026-03-12T10:00:00Z",
      "updated_at": "2026-03-12T10:00:00Z"
    }
  ]
}
```

#### Get Product by ID
```http
GET /api/products/get?id={uuid}
```

### gRPC API

**Address**: `localhost:50051`

#### Service Definition
```protobuf
service ProductService {
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
}
```

## Setup Instructions

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Protocol Buffers compiler (`protoc`)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install protoc plugins (if not installed)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Database Setup

1. Create PostgreSQL database:
```sql
CREATE DATABASE vave_db;
```

2. Run migrations:
```bash
make migrate-up
```

Or manually:
```bash
psql -h localhost -U postgres -d vave_db -f migrations/000003_create_products.up.sql
```

### Environment Configuration

1. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

2. Update database credentials if needed

### Generate Protobuf Code

```bash
chmod +x scripts/generate_proto.sh
./scripts/generate_proto.sh
```

Or use Make:
```bash
make proto
```

### Run the Application

```bash
# Using Make
make run

# Or directly
go run cmd/api/main.go
```

The application will start:
- REST API on `http://localhost:8080`
- gRPC server on `localhost:50051`

## Testing

### Test REST API

```bash
# List products
curl http://localhost:8080/api/products

# Health check
curl http://localhost:8080/health
```

### Test gRPC API

Using `grpcurl`:

```bash
# List services
grpcurl -plaintext localhost:50051 list

# List products
grpcurl -plaintext localhost:50051 product.ProductService/ListProducts
```

## Project Structure

```
backend/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── domain/                  # Domain models & interfaces
│   │   ├── product.go
│   │   └── errors.go
│   ├── repository/              # Data access layer
│   │   └── product.go
│   ├── service/                 # Business logic layer
│   │   └── product.go
│   ├── api/                     # REST API
│   │   ├── handler/
│   │   │   └── product.go
│   │   ├── router/
│   │   │   └── router.go
│   │   └── response/
│   │       └── response.go
│   ├── grpc/                    # gRPC handlers
│   │   ├── product_server.go
│   │   └── pb/                  # Generated protobuf code
│   ├── config/                  # Configuration
│   │   └── config.go
│   └── pkg/                     # Internal packages
│       └── db/
│           └── postgres.go
├── proto/                       # Protobuf definitions
│   └── product.proto
├── migrations/                  # Database migrations
│   └── 000003_create_products.*.sql
└── scripts/                     # Helper scripts
    └── generate_proto.sh
```

## Development

### Add Sample Data

```sql
INSERT INTO products (name, description, price, stock_quantity, category, sku)
VALUES 
  ('Laptop', 'High-performance laptop', 1299.99, 50, 'Electronics', 'LAP-001'),
  ('Mouse', 'Wireless gaming mouse', 59.99, 200, 'Electronics', 'MOU-001'),
  ('Keyboard', 'Mechanical keyboard', 129.99, 150, 'Electronics', 'KEY-001');
```

### Make Commands

```bash
make build        # Build the application
make run          # Run the application
make test         # Run tests
make proto        # Generate protobuf code
make migrate-up   # Run migrations
make migrate-down # Rollback migrations
make deps         # Download dependencies
make clean        # Clean build artifacts
```

## Architecture

The application follows Clean Architecture principles:

1. **Domain Layer**: Core business entities and interfaces
2. **Repository Layer**: Database operations (SQL queries)
3. **Service Layer**: Business logic orchestration
4. **Handler Layer**: HTTP/gRPC request handling
5. **Router**: API endpoint registration

## License

MIT
