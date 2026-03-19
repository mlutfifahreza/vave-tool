# Internal API Authentication

## Overview

The Vave Tool API now has database-backed authentication for internal write operations (POST, PUT, DELETE). Client credentials are stored in a `clients` table and all write operations are tracked with the client ID.

## Setup

### 1. Run Migrations

```bash
make migrate-up
```

This will create:
- `clients` table for storing API client credentials
- `updated_by` field in the `products` table

### 2. Seed Admin Client

```bash
make seed-client
```

This creates/updates an admin client with:
- **Username**: `admin`
- **Password**: `admin123`

You can modify the credentials in `script/seed_client.go` before running the seed command.

## API Endpoints

### Public Endpoints (No Authentication)
- `GET /api/products` - List all products
- `GET /api/products/{id}` - Get product by ID

### Internal Endpoints (Basic Auth Required)
- `POST /internal/products` - Create product
- `PUT /internal/products/{id}` - Update product
- `DELETE /internal/products/{id}` - Delete product

## Using Internal Endpoints

### cURL Example

```bash
# Create a product
curl -X POST http://localhost:8080/internal/products \
  -u admin:admin123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "description": "Product description",
    "price": 99.99,
    "stock_quantity": 100,
    "category": "Electronics",
    "sku": "PROD-001",
    "is_active": true
  }'

# Update a product
curl -X PUT http://localhost:8080/internal/products/{id} \
  -u admin:admin123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "price": 89.99
  }'

# Delete a product
curl -X DELETE http://localhost:8080/internal/products/{id} \
  -u admin:admin123
```

### Python Example

```python
import requests
from requests.auth import HTTPBasicAuth

auth = HTTPBasicAuth('admin', 'admin123')

# Create product
response = requests.post(
    'http://localhost:8080/internal/products',
    auth=auth,
    json={
        'name': 'New Product',
        'price': 99.99,
        'stock_quantity': 100,
        'is_active': True
    }
)

print(response.json())
```

## Client Tracking

Every write operation (create, update, delete) is tracked with the authenticated client's ID in the `updated_by` field of the product record. This provides an audit trail of who made changes to products.

## Security Features

- Passwords are hashed using bcrypt
- Constant-time comparison prevents timing attacks
- Only active clients can authenticate
- Client lookup happens on every request (no session storage)

## Database Schema

### Clients Table

```sql
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Products Table (Updated)

```sql
ALTER TABLE products ADD COLUMN updated_by UUID;
ALTER TABLE products ADD CONSTRAINT fk_products_updated_by 
    FOREIGN KEY (updated_by) REFERENCES clients(id) ON DELETE SET NULL;
```

## Adding New Clients

You can manually insert new clients using the database or create a simple script similar to `script/seed_client.go`:

```bash
# Example: Create a new client directly in the database
psql -d vave_db -c "
INSERT INTO clients (name, username, password, is_active)
VALUES ('Client Name', 'username', '\$bcrypt_hash', true);
"
```

To generate a bcrypt hash in Go:

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "your_password"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```
