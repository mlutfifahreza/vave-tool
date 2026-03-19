# 1. List all products (with pagination)
curl -X GET "http://localhost:8080/api/products?page=1&size=10"

# 2. Get product by ID
curl -X GET "http://localhost:8080/api/products/{product-id-here}"

# 3. Create a new product (requires authentication)
curl -X POST http://localhost:8080/internal/products \
  -u admin:admin123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Headphones",
    "description": "High-quality wireless headphones with noise cancellation",
    "price": 199.99,
    "stock_quantity": 50,
    "category": "Electronics",
    "sku": "WH-2024-001",
    "is_active": true
  }'

# 4. Update a product (requires authentication)
curl -X PUT http://localhost:8080/internal/products/{product-id-here} \
  -u admin:admin123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Headphones Pro",
    "description": "Premium wireless headphones with advanced noise cancellation",
    "price": 249.99,
    "stock_quantity": 30,
    "category": "Electronics",
    "sku": "WH-2024-001",
    "is_active": true
  }'

# 5. Delete a product (requires authentication)
curl -X DELETE http://localhost:8080/internal/products/{product-id-here} \
  -u admin:admin123

# Additional examples:

# Health check
curl http://localhost:8080/health

# Metrics endpoint
curl http://localhost:8080/metrics

# Create product with minimal fields
curl -X POST http://localhost:8080/internal/products \
  -u admin:admin123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Simple Product",
    "price": 29.99,
    "stock_quantity": 100,
    "is_active": true
  }'

# Unauthorized request (will get 401)
curl -X POST http://localhost:8080/internal/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "price": 10, "stock_quantity": 1}'