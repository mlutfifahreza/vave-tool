# 1. List all products (with pagination)
curl -X GET "http://localhost:8080/api/products?page=1&size=10"

<!-- response: -->
{
  "success": true,
  "data": {
    "products": [
      {
        "id": "818ba0f4-5d08-473e-b0bb-b9067196de09",
        "name": "Portable Water Bottle 999858",
        "description": "Insulated water bottle, 32oz capacity",
        "price": 41.97,
        "stock_quantity": 109,
        "category": "Accessories",
        "sku": "SKU-ACC-999858",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "7a066e96-f246-44dc-b85e-1eb3f8805067",
        "name": "Standard Laptop Pro 999859",
        "description": "High-performance laptop with 8GB RAM and 16GB SSD",
        "price": 1593.07,
        "stock_quantity": 423,
        "category": "Electronics",
        "sku": "SKU-ELE-999859",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "ebfa9b28-5505-4384-9616-ca2518afe2ff",
        "name": "Premium Desk Lamp 999861",
        "description": "16 desk lamp with adjustable brightness",
        "price": 128.77,
        "stock_quantity": 205,
        "category": "Furniture",
        "sku": "SKU-FUR-999861",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "dd02c078-9e5e-469b-972f-082a3cd552c8",
        "name": "Blender 999862",
        "description": "24-speed blender with 48 cups capacity",
        "price": 139.47,
        "stock_quantity": 461,
        "category": "Appliances",
        "sku": "SKU-APP-999862",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "2ad2fc29-3c3c-469d-8549-1008912044f4",
        "name": "Notebook Set 999863",
        "description": "Set of 48 ruled notebooks, 200 pages each",
        "price": 11.66,
        "stock_quantity": 496,
        "category": "Stationery",
        "sku": "SKU-STA-999863",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "db8dc8e6-99c9-49ae-96ab-3a67b76af64c",
        "name": "Bluetooth Speaker 999864",
        "description": "Portable Bluetooth speaker with 32-hour battery",
        "price": 134.04,
        "stock_quantity": 349,
        "category": "Electronics",
        "sku": "SKU-ELE-999864",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "c4e74fdc-3816-4005-bb8d-fa99a12690f8",
        "name": "Notebook Set 999866",
        "description": "Set of 12 ruled notebooks, 200 pages each",
        "price": 14.31,
        "stock_quantity": 145,
        "category": "Stationery",
        "sku": "SKU-STA-999866",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "1feb8618-e727-41bc-a5c1-82331c8120bf",
        "name": "Standard Desk Lamp 999867",
        "description": "16 desk lamp with adjustable brightness",
        "price": 46.37,
        "stock_quantity": 177,
        "category": "Furniture",
        "sku": "SKU-FUR-999867",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "ac15eb5d-33eb-46e3-b2c7-cbc25876646e",
        "name": "Microwave 999868",
        "description": "12 watt microwave oven with 24 presets",
        "price": 91.81,
        "stock_quantity": 266,
        "category": "Appliances",
        "sku": "SKU-APP-999868",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      },
      {
        "id": "f5ede528-86e9-4f55-bda4-221800d4ec96",
        "name": "Professional Bookshelf 999869",
        "description": "24-tier 48 bookshelf",
        "price": 102.16,
        "stock_quantity": 252,
        "category": "Furniture",
        "sku": "SKU-FUR-999869",
        "is_active": true,
        "created_at": "2026-03-12T23:27:42.48478Z",
        "updated_at": "2026-03-12T23:27:42.48478Z"
      }
    ]
  }
}

# 2. Get product by ID
curl -X GET "http://localhost:8080/api/products/{product-id-here}"

<!-- response: -->
{
  "success": true,
  "data": {
    "id": "ac15eb5d-33eb-46e3-b2c7-cbc25876646e",
    "name": "Microwave 999868",
    "description": "12 watt microwave oven with 24 presets",
    "price": 91.81,
    "stock_quantity": 266,
    "category": "Appliances",
    "sku": "SKU-APP-999868",
    "is_active": true,
    "created_at": "2026-03-12T23:27:42.48478Z",
    "updated_at": "2026-03-12T23:27:42.48478Z"
  }
}

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

<!-- response: -->
{
  "success": true,
  "data": {
    "id": "9a27a12e-7c02-4994-bdb6-7297175b0087",
    "name": "Wireless Headphones",
    "description": "High-quality wireless headphones with noise cancellation",
    "price": 199.99,
    "stock_quantity": 50,
    "category": "Electronics",
    "sku": "WH-2024-001",
    "is_active": true,
    "created_at": "2026-03-20T02:43:48.692772Z",
    "updated_at": "2026-03-20T02:43:48.692772Z"
  }
}

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

<!-- response: -->
{
  "success": true,
  "data": {
    "id": "9a27a12e-7c02-4994-bdb6-7297175b0087",
    "name": "Wireless Headphones Pro",
    "description": "Premium wireless headphones with advanced noise cancellation",
    "price": 249.99,
    "stock_quantity": 30,
    "category": "Electronics",
    "sku": "WH-2024-001",
    "is_active": true,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2026-03-20T02:44:25.204404Z"
  }
}

# 5. Delete a product (requires authentication)
curl -X DELETE http://localhost:8080/internal/products/{product-id-here} \
  -u admin:admin123

<!-- response: -->
{
  "success": true,
  "data": {
    "success": true
  }
}