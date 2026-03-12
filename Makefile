.PHONY: help build run test proto migrate-up migrate-down seed-products swagger

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make proto       - Generate protobuf files"
	@echo "  make swagger     - Generate Swagger documentation"
	@echo "  make migrate-up  - Run database migrations up"
	@echo "  make migrate-down - Run database migrations down"
	@echo "  make seed-products - Insert dummy products (dev only)"

run:
	go run cmd/api/main.go

test:
	go test -v ./...

proto:
	protoc --go_out=. --go_opt=module=github.com/vave-tool/backend \
		--go-grpc_out=. --go-grpc_opt=module=github.com/vave-tool/backend \
		proto/product.proto

swagger:
	swag init -g cmd/api/main.go -o docs

migrate-up:
	@echo "Running migrations..."
	@for file in migrations/*up.sql; do \
		echo "Applying $$file"; \
		psql -h localhost -U postgres -d vave_db -f $$file; \
	done

migrate-down:
	@echo "Rolling back migrations..."
	@for file in $$(ls -r migrations/*down.sql); do \
		echo "Applying $$file"; \
		psql -h localhost -U postgres -d vave_db -f $$file; \
	done

seed-products:
	@echo "Inserting dummy products..."
	@psql -h localhost -U postgres -d vave_db -c "\
		INSERT INTO products (name, description, price, stock_quantity, category, sku, is_active) VALUES \
		('Laptop Pro 15', 'High-performance laptop with 16GB RAM and 512GB SSD', 1299.99, 25, 'Electronics', 'LAP-PRO-15', true), \
		('Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 29.99, 150, 'Electronics', 'MOUSE-WL-01', true), \
		('Office Chair', 'Adjustable ergonomic office chair with lumbar support', 349.50, 45, 'Furniture', 'CHAIR-OFF-01', true), \
		('Coffee Maker', 'Programmable 12-cup coffee maker with timer', 89.99, 60, 'Appliances', 'COFFEE-12C', true), \
		('Notebook Set', 'Set of 3 ruled notebooks, 200 pages each', 12.99, 200, 'Stationery', 'NOTE-SET-3', true), \
		('USB-C Cable', '2-meter USB-C to USB-C charging cable', 19.99, 300, 'Electronics', 'CABLE-USBC-2M', true), \
		('Desk Lamp', 'LED desk lamp with adjustable brightness', 45.00, 80, 'Furniture', 'LAMP-LED-ADJ', true), \
		('Water Bottle', 'Stainless steel insulated water bottle, 24oz', 24.99, 120, 'Accessories', 'BOTTLE-SS-24', true), \
		('Bluetooth Speaker', 'Portable Bluetooth speaker with 10-hour battery', 79.99, 55, 'Electronics', 'SPEAK-BT-10H', true), \
		('Standing Desk', 'Electric height-adjustable standing desk', 599.00, 15, 'Furniture', 'DESK-STAND-EL', true) \
		ON CONFLICT (sku) DO NOTHING;"
	@echo "Dummy products inserted successfully!"

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/
