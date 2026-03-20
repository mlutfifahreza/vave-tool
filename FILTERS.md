# Product List Filters

The ListProducts endpoint has been enhanced with filter capabilities for category, subcategory, and price range.

## API Endpoint

```
GET /api/products
```

## Query Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `page` | int | Page number (default: 1) | `page=1` |
| `size` | int | Items per page (default: 10, max: 100) | `size=20` |
| `category_id` | string | Filter by category UUID | `category_id=123e4567-e89b-12d3-a456-426614174000` |
| `subcategory_id` | string | Filter by subcategory UUID | `subcategory_id=987fcdeb-51a2-43d7-b234-123456789abc` |
| `min_price` | float | Minimum price filter | `min_price=10.00` |
| `max_price` | float | Maximum price filter | `max_price=100.00` |

## Examples

### 1. List all products (paginated)
```bash
curl "http://localhost:8080/api/products?page=1&size=10"
```

### 2. Filter by category
```bash
curl "http://localhost:8080/api/products?category_id=123e4567-e89b-12d3-a456-426614174000"
```

### 3. Filter by subcategory
```bash
curl "http://localhost:8080/api/products?subcategory_id=987fcdeb-51a2-43d7-b234-123456789abc"
```

### 4. Filter by price range
```bash
curl "http://localhost:8080/api/products?min_price=50.00&max_price=200.00"
```

### 5. Combine multiple filters
```bash
curl "http://localhost:8080/api/products?category_id=123e4567-e89b-12d3-a456-426614174000&min_price=10.00&max_price=100.00&page=1&size=20"
```

### 6. Filter by category and price range
```bash
curl "http://localhost:8080/api/products?category_id=123e4567-e89b-12d3-a456-426614174000&min_price=25.00&max_price=75.00"
```

## Implementation Details

### Cache Key Strategy
The cache key now includes all filter parameters to ensure proper cache isolation:

```
products:list:page:1:size:10:cat:{category_id}:subcat:{subcategory_id}:minp:{min_price}:maxp:{max_price}
```

### SQL Query Optimization
The repository builds dynamic SQL queries based on the provided filters:

```sql
SELECT p.id, p.name, p.description, p.price, p.stock_quantity, 
       p.category_id, c.name as category_name,
       p.subcategory_id, s.name as subcategory_name,
       p.sku, p.is_active, p.updated_by, p.created_at, p.updated_at
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN subcategories s ON p.subcategory_id = s.id
WHERE p.is_active = true
  AND p.category_id = $1        -- Only if category_id is provided
  AND p.subcategory_id = $2     -- Only if subcategory_id is provided
  AND p.price >= $3             -- Only if min_price is provided
  AND p.price <= $4             -- Only if max_price is provided
ORDER BY p.created_at DESC
LIMIT $5 OFFSET $6
```

## Changes Made

1. **Domain Layer** ([internal/domain/product.go](internal/domain/product.go))
   - Added `ProductFilterParams` struct
   - Updated `ProductRepository.List()` and `ProductService.ListProducts()` signatures

2. **Repository Layer** ([internal/repository/product.go](internal/repository/product.go))
   - Enhanced `List()` method to build dynamic SQL with filter conditions
   - Updated `Count()` method to support filtered counting

3. **Service Layer** ([internal/service/product.go](internal/service/product.go))
   - Added `generateCacheKey()` helper to create filter-aware cache keys
   - Updated `ListProducts()` to accept and pass filters

4. **Handler Layer** ([internal/api/handler/product.go](internal/api/handler/product.go))
   - Added query parameter parsing for all filter options
   - Enhanced Swagger documentation

5. **gRPC Server** ([internal/grpc/product_server.go](internal/grpc/product_server.go))
   - Updated to pass empty filters for backward compatibility

## Performance Considerations

- Filters are applied at the database level for optimal performance
- Cache keys are unique per filter combination to prevent cache collisions
- Indexes on `category_id`, `subcategory_id`, and `price` columns are recommended for better query performance
